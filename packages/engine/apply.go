package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

func (e *Engine) Apply(event event.GameEvent) error {
	e.log.Debug("Applying event:", event.EventType())
	e.record.Add(event)
	newGame, err := e.applyEvent(e.game, event)
	if err != nil {
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
	switch evnt := gameEvent.(type) {
	case event.GameLifecycleEvent:
		return e.applyGameLifecycleEvent(game, evnt)
	case event.PriorityEvent:
		return e.applyPriorityEvent(game, evnt)
	case event.TurnEvent:
		return e.applyTurnEvent(game, evnt)
	case event.DrawCardEvent:
		return ApplyDrawCardEvent(game, evnt)
	case event.UntapAllEvent:
		return ApplyUntapAllEvent(game, evnt)
	case event.SetNextPlayerEvent:
		return ApplySetNextPlayerEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown event type: %T", gameEvent)
	}
	return game, nil
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
		return game, nil
	case event.GameOverEvent:
		e.log.Info("Game over, winner:", evnt.WinnerID)
		newGame := game.WithGameOver(evnt.WinnerID)
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown game lifecycle event type: %T", evnt)
	}
	return game, fmt.Errorf("applyGameLifecycleEvent: unknown game lifecycle event type: %T", gameLifecycleEvent)
}

func (e *Engine) applyTurnEvent(game state.Game, turnEvent event.TurnEvent) (state.Game, error) {
	switch evnt := turnEvent.(type) {
	case event.BeginStepEvent:
		return e.applyBeginStepEvent(game, evnt)
	case event.EndStepEvent:
		return e.applyEndStepEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown turn event type: %T", evnt)
	}
	return game, fmt.Errorf("applyTurnEvent: unknown turn event type: %T", turnEvent)
}

func (e *Engine) applyPriorityEvent(game state.Game, priorityEvent event.PriorityEvent) (state.Game, error) {
	switch evnt := priorityEvent.(type) {
	case event.AllPlayersPassedPriorityEvent:
		return game, nil
	case event.PassPriorityEvent:
		// TODO Get PlayerID from one source - event vs state
		fmt.Println("Player with priority:", game.PriorityPlayerID())
		nextPlayerIDWithPriority := game.NextPlayerID(game.PriorityPlayerID())
		fmt.Println("Next player with priority:", nextPlayerIDWithPriority)
		newGame := game.WithPlayerPassedPriority(
			evnt.PlayerID,
		).WithPlayerWithPriority(
			nextPlayerIDWithPriority,
		)
		fmt.Println("New Player with priority:", newGame.PriorityPlayerID())
		return newGame, nil
		return game.WithPlayerPassedPriority(evnt.PlayerID), nil
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
	newGame := game.WithClearedPriority().WithResetPriorityPasses()
	return newGame, nil
}

func (e *Engine) applyBeginStepEvent(
	game state.Game,
	beginStepEvent event.BeginStepEvent,
) (state.Game, error) {
	switch beginStepEvent.(type) {
	case event.BeginUntapStepEvent:
		return game, nil
	case event.BeginUpkeepStepEvent:
		return game, nil
	case event.BeginDrawStepEvent:
		return game, nil
	case event.BeginPrecombatMainStepEvent:
		return game, nil
	case event.BeginBeginningOfCombatStepEvent:
		return game, nil
	case event.BeginDeclareAttackersStepEvent:
		return game, nil
	case event.BeginDeclareBlockersStepEvent:
		return game, nil
	case event.BeginCombatDamageStepEvent:
		return game, nil
	case event.BeginEndOfCombatStepEvent:
		return game, nil
	case event.BeginEndStepEvent:
		return game, nil
	case event.BeginCleanupStepEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown begin step event type: %T", beginStepEvent)
	}
	return game, nil
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
	if err != nil {
		return game, fmt.Errorf("untap: %w", err)
	}
	newGame := game.WithBattlefield(newBattlefield)
	return newGame, nil
}

func ApplyDrawCardEvent(
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
