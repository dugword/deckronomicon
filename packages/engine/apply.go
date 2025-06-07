package engine

// TODO: Document what level things should live at. Maybe apply is where the
// core game engine logic and enforcement lives. it takes the structured
// imput, verifies per the rules of the game it can happen, and then applies
// it.

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func (e *Engine) Apply(event event.GameEvent) error {
	e.log.Info("Applying event:", event.EventType())
	newGame, err := e.applyEvent(e.game, event)
	if err != nil {
		if errors.Is(err, mtg.ErrGameOver) {
			e.log.Info("Game over detected")
			return err
		}
		e.log.Critical("Failed to apply event:", err)
		return fmt.Errorf("apply event: %w", err)
	}
	e.game = newGame
	for _, agent := range e.agents {
		if err := agent.ReportState(e.game); err != nil {
			e.log.Critical("Failed to report state to agent:", err)
			return fmt.Errorf("report state: %w", err)
		}
	}
	return nil
}

func (e *Engine) applyEvent(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
	e.record.Add(gameEvent)
	switch evnt := gameEvent.(type) {
	case event.GameLifecycleEvent:
		return e.applyGameLifecycleEvent(game, evnt)
	case event.PriorityEvent:
		return e.applyPriorityEvent(game, evnt)
	case event.TurnEvent:
		return e.applyTurnEvent(game, evnt)
	case event.TurnBasedActionEvent:
		return e.applyTurnBasedActionEvent(game, evnt)
	case event.CombatEvent:
		return e.applyCombatEvent(game, evnt)

	case event.SetNextPlayerEvent:
		return ApplySetNextPlayerEvent(game, evnt)
	default:
		return e.applyOtherEvents(game, evnt)
	}
}

func (e *Engine) applyTurnBasedActionEvent(game state.Game, turnBasedActionEvent event.TurnBasedActionEvent) (state.Game, error) {
	switch evnt := turnBasedActionEvent.(type) {
	case event.UntapAllEvent:
		return ApplyUntapAllEvent(game, evnt)
	case event.UpkeepEvent:
		return game, nil
	case event.ProgressSagaEvent:
		return game, nil
	case event.CheckDayNightEvent:
		return game, nil
	case event.PhaseInPhaseOutEvent:
		return game, nil
	case event.DiscardToHandSizeEvent:
		return game, nil
	case event.RemoveDamageEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown turn-based action event type: %T", evnt)
	}
}

func (e *Engine) applyOtherEvents(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
	switch evnt := gameEvent.(type) {
	case event.ConcedeEvent:
		return e.ApplyConcedeEvent(game, evnt)
	case event.PlayLandEvent:
		return e.applyPlayLandEvent(game, evnt)
	case event.CastSpellEvent:
		return e.ApplyCastSpellEvent(game, evnt)
	case event.DrawCardEvent:
		return e.ApplyDrawCardEvent(game, evnt)
	case event.DiscardCardEvent:
		return e.ApplyDiscardCardEvent(game, evnt)
	case event.DrawStartingHandEvent:
		return game, nil
	case event.ShuffleDeckEvent:
		player, err := game.GetPlayer(evnt.PlayerID)
		if err != nil {
			return game, fmt.Errorf("shuffle deck: %w", err)
		}
		newPlayer := player.WithShuffleDeck(e.rng.DeckShuffler())
		newGame := game.WithUpdatedPlayer(newPlayer)
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown event type: %T", evnt)
	}
}

func (e *Engine) applyCombatEvent(game state.Game, combatEvent event.CombatEvent) (state.Game, error) {
	switch evnt := combatEvent.(type) {
	case event.DeclareAttackersEvent:
		return game, nil
	case event.DeclareBlockersEvent:
		return game, nil
	case event.CombatDamageEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown combat event type: %T", evnt)
	}
}

func (e *Engine) applyGameLifecycleEvent(game state.Game, gameLifecycleEvent event.GameLifecycleEvent) (state.Game, error) {
	switch evnt := gameLifecycleEvent.(type) {
	case event.BeginGameEvent:
		e.log.Info("Game started")
		return game, nil
	case event.BeginTurnEvent:
		player, err := game.GetPlayer(evnt.PlayerID)
		if err != nil {
			return game, fmt.Errorf("begin turn: %w", err)
		}
		newPlayer := player.WithNextTurn()
		newGame := game.WithUpdatedPlayer(newPlayer)
		return newGame, nil
	case event.EndTurnEvent:
		return game, nil
	case event.GameOverEvent:
		e.log.Info("Game over, winner:", evnt.WinnerID)
		newGame := game.WithGameOver(evnt.WinnerID)
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown game lifecycle event type: %T", evnt)
	}
}

func (e *Engine) applyTurnEvent(game state.Game, turnEvent event.TurnEvent) (state.Game, error) {
	switch evnt := turnEvent.(type) {
	case event.BeginPhaseEvent:
		return e.applyBeginPhaseEvent(game, evnt)
	case event.EndPhaseEvent:
		return e.applyEndPhaseEvent(game, evnt)
	case event.BeginStepEvent:
		return e.applyBeginStepEvent(game, evnt)
	case event.EndStepEvent:
		return e.applyEndStepEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown turn event type: %T", evnt)
	}
}

func (e *Engine) applyPriorityEvent(game state.Game, priorityEvent event.PriorityEvent) (state.Game, error) {
	switch evnt := priorityEvent.(type) {
	case event.AllPlayersPassedPriorityEvent:
		return game, nil
	case event.PassPriorityEvent:
		// TODO Get PlayerID from one source - event vs state
		nextPlayerIDWithPriority := game.NextPlayerID(game.PriorityPlayerID())
		newGame := game.WithPlayerPassedPriority(
			evnt.PlayerID,
		).WithPlayerWithPriority(
			nextPlayerIDWithPriority,
		)
		return newGame, nil
	case event.ReceivePriorityEvent:
		newGame := game.WithPlayerWithPriority(
			evnt.PlayerID,
		)
		return newGame, nil
	case event.ResetPriorityPassesEvent:
		newGame := game.WithResetPriorityPasses()
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown priority event type: %T", evnt)
	}
}

func (e *Engine) applyEndStepEvent(
	game state.Game,
	endStepEvent event.EndStepEvent,
) (state.Game, error) {
	game = game.WithClearedPriority().WithResetPriorityPasses()
	switch endStepEvent.(type) {
	case event.EndUntapStepEvent:
		return game, nil
	case event.EndUpkeepStepEvent:
		return game, nil
	case event.EndDrawStepEvent:
		return game, nil
	case event.EndPrecombatMainStepEvent:
		return game, nil
	case event.EndBeginningOfCombatStepEvent:
		return game, nil
	case event.EndDeclareAttackersStepEvent:
		return game, nil
	case event.EndDeclareBlockersStepEvent:
		return game, nil
	case event.EndCombatDamageStepEvent:
		return game, nil
	case event.EndEndOfCombatStepEvent:
		return game, nil
	case event.EndPostcombatMainStepEvent:
		return game, nil
	case event.EndEndStepEvent:
		return game, nil
	case event.EndCleanupStepEvent:
		return game, nil
	}
	return game, nil
}

func (e *Engine) applyBeginStepEvent(
	game state.Game,
	beginStepEvent event.BeginStepEvent,
) (state.Game, error) {
	switch beginStepEvent.(type) {
	case event.BeginUntapStepEvent:
		newGame := game.WithStep(mtg.StepUntap)
		return newGame, nil
	case event.BeginUpkeepStepEvent:
		newGame := game.WithStep(mtg.StepUpkeep)
		return newGame, nil
	case event.BeginDrawStepEvent:
		newGame := game.WithStep(mtg.StepDraw)
		return newGame, nil
	case event.BeginPrecombatMainStepEvent:
		newGame := game.WithStep(mtg.StepPrecombatMain)
		return newGame, nil
	case event.BeginBeginningOfCombatStepEvent:
		newGame := game.WithStep(mtg.StepBeginningOfCombat)
		return newGame, nil
	case event.BeginDeclareAttackersStepEvent:
		newGame := game.WithStep(mtg.StepDeclareAttackers)
		return newGame, nil
	case event.BeginDeclareBlockersStepEvent:
		newGame := game.WithStep(mtg.StepDeclareBlockers)
		return newGame, nil
	case event.BeginCombatDamageStepEvent:
		newGame := game.WithStep(mtg.StepCombatDamage)
		return newGame, nil
	case event.BeginEndOfCombatStepEvent:
		newGame := game.WithStep(mtg.StepEndOfCombat)
		return newGame, nil
	case event.BeginPostcombatMainStepEvent:
		newGame := game.WithStep(mtg.StepPostcombatMain)
		return newGame, nil
	case event.BeginEndStepEvent:
		newGame := game.WithStep(mtg.StepEnd)
		return newGame, nil
	case event.BeginCleanupStepEvent:
		newGame := game.WithStep(mtg.StepCleanup)
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown begin step event type: %T", beginStepEvent)
	}
}

func (e *Engine) applyEndPhaseEvent(
	game state.Game,
	endPhaseEvent event.EndPhaseEvent,
) (state.Game, error) {
	switch endPhaseEvent.(type) {
	case event.EndBeginningPhaseEvent:
		return game, nil
	case event.EndPrecombatMainPhaseEvent:
		return game, nil
	case event.EndCombatPhaseEvent:
		return game, nil
	case event.EndPostcombatMainPhaseEvent:
		return game, nil
	case event.EndEndingPhaseEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown end phase event type: %T", endPhaseEvent)
	}
}
func (e *Engine) applyBeginPhaseEvent(
	game state.Game,
	beginPhaseEvent event.BeginPhaseEvent,
) (state.Game, error) {
	switch beginPhaseEvent.(type) {
	case event.BeginBeginningPhaseEvent:
		newGame := game.WithPhase(mtg.PhaseBeginning)
		return newGame, nil
	case event.BeginPrecombatMainPhaseEvent:
		newGame := game.WithPhase(mtg.PhasePrecombatMain)
		return newGame, nil
	case event.BeginCombatPhaseEvent:
		newGame := game.WithPhase(mtg.PhaseCombat)
		return newGame, nil
	case event.BeginPostcombatMainPhaseEvent:
		newGame := game.WithPhase(mtg.PhasePostcombatMain)
		return newGame, nil
	case event.BeginEndingPhaseEvent:
		newGame := game.WithPhase(mtg.PhaseEnding)
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown begin phase event type: %T", beginPhaseEvent)
	}
}

func ApplySetNextPlayerEvent(
	game state.Game,
	event event.SetNextPlayerEvent,
) (state.Game, error) {
	nextPlayer := game.NextPlayerID(game.ActivePlayerID())
	newGame := game.WithActivePlayer(nextPlayer)
	return newGame, nil
}

func ApplyUntapAllEvent(
	game state.Game,
	event event.UntapAllEvent,
) (state.Game, error) {
	player, err := game.GetPlayer(event.PlayerID)
	if err != nil {
		return game, fmt.Errorf("untap: %w", err)
	}
	newBattlefield := game.Battlefield()
	newBattlefield.UntapAll(player.ID())
	newGame := game.WithBattlefield(newBattlefield)
	return newGame, nil
}

func (e *Engine) ApplyDrawCardEvent(
	game state.Game,
	event event.DrawCardEvent,
) (state.Game, error) {
	player, err := game.GetPlayer(event.PlayerID)
	if err != nil {
		return game, fmt.Errorf("draw: %w", err)
	}
	player, _, err = player.WithDrawCard()
	if err != nil {
		return game, fmt.Errorf("draw: %w", err)
	}
	newGame := game.WithUpdatedPlayer(player)
	return newGame, nil
}

func (e *Engine) ApplyDiscardCardEvent(
	game state.Game,
	event event.DiscardCardEvent,
) (state.Game, error) {
	player, err := game.GetPlayer(event.PlayerID)
	if err != nil {
		return game, fmt.Errorf("discard: %w", err)
	}
	player, err = player.WithDiscardCard(event.CardID)
	if err != nil {
		return game, fmt.Errorf("discard: %w", err)
	}
	newGame := game.WithUpdatedPlayer(player)
	return newGame, nil
}

func (e *Engine) ApplyCastSpellEvent(
	game state.Game,
	event event.CastSpellEvent,
) (state.Game, error) {
	player, err := game.GetPlayer(event.PlayerID)
	if err != nil {
		return game, fmt.Errorf("cast spell: %w", err)
	}
	/*
		player, err = player.WithCardOnStack(event.CardID)
		if err != nil {
			return game, fmt.Errorf("cast spell: %w", err)
		}
	*/
	newGame := game.WithUpdatedPlayer(player)
	return newGame, nil
}

func (e *Engine) ApplyConcedeEvent(
	game state.Game,
	event event.ConcedeEvent,
) (state.Game, error) {
	opponent, err := game.GetOpponent(event.PlayerID)
	if err != nil {
		return game, fmt.Errorf("concede: %w", err)
	}
	newGame := game.WithGameOver(opponent.ID())
	return newGame, mtg.ErrGameOver
}
