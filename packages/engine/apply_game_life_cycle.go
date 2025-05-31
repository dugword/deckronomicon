package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

func (e *Engine) applyGameLifecycleEvent(game state.Game, gameLifecycleEvent event.GameLifecycleEvent) (state.Game, error) {
	switch evnt := gameLifecycleEvent.(type) {
	case event.TurnEvent:
		return e.applyTurnEvent(game, evnt)
	case event.BeginGameEvent:
		return e.applyBeginGameEvent(game, evnt)
	case event.EndGameEvent:
		return e.applyEndGameEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown game lifecycle event '%T'", evnt)
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
	case event.BeginTurnEvent:
		return e.applyBeginTurnEvent(game, evnt)
	case event.EndTurnEvent:
		return e.applyEndTurnEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown turn event '%T'", evnt)
	}
}

func (e *Engine) applyBeginPhaseEvent(
	game state.Game,
	beginPhaseEvent event.BeginPhaseEvent,
) (state.Game, error) {
	switch beginPhaseEvent.(type) {
	case event.BeginBeginningPhaseEvent:
		game = game.WithPhase(mtg.PhaseBeginning)
		return game, nil
	case event.BeginPrecombatMainPhaseEvent:
		game = game.WithPhase(mtg.PhasePrecombatMain)
		return game, nil
	case event.BeginCombatPhaseEvent:
		game = game.WithPhase(mtg.PhaseCombat)
		return game, nil
	case event.BeginPostcombatMainPhaseEvent:
		game = game.WithPhase(mtg.PhasePostcombatMain)
		return game, nil
	case event.BeginEndingPhaseEvent:
		game = game.WithPhase(mtg.PhaseEnding)
		return game, nil
	default:
		return game, fmt.Errorf("unknown begin phase event '%T'", beginPhaseEvent)
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
		return game, fmt.Errorf("unknown end phase event '%T'", endPhaseEvent)
	}
}

func (e *Engine) applyBeginStepEvent(
	game state.Game,
	beginStepEvent event.BeginStepEvent,
) (state.Game, error) {
	switch beginStepEvent.(type) {
	case event.BeginUntapStepEvent:
		game = game.WithStep(mtg.StepUntap)
		return game, nil
	case event.BeginUpkeepStepEvent:
		game = game.WithStep(mtg.StepUpkeep)
		return game, nil
	case event.BeginDrawStepEvent:
		game = game.WithStep(mtg.StepDraw)
		return game, nil
	case event.BeginPrecombatMainStepEvent:
		game = game.WithStep(mtg.StepPrecombatMain)
		return game, nil
	case event.BeginBeginningOfCombatStepEvent:
		game = game.WithStep(mtg.StepBeginningOfCombat)
		return game, nil
	case event.BeginDeclareAttackersStepEvent:
		game = game.WithStep(mtg.StepDeclareAttackers)
		return game, nil
	case event.BeginDeclareBlockersStepEvent:
		game = game.WithStep(mtg.StepDeclareBlockers)
		return game, nil
	case event.BeginCombatDamageStepEvent:
		game = game.WithStep(mtg.StepCombatDamage)
		return game, nil
	case event.BeginEndOfCombatStepEvent:
		game = game.WithStep(mtg.StepEndOfCombat)
		return game, nil
	case event.BeginPostcombatMainStepEvent:
		game = game.WithStep(mtg.StepPostcombatMain)
		return game, nil
	case event.BeginEndStepEvent:
		game = game.WithStep(mtg.StepEnd)
		return game, nil
	case event.BeginCleanupStepEvent:
		game = game.WithStep(mtg.StepCleanup)
		return game, nil
	default:
		return game, fmt.Errorf("unknown begin step event '%T'", beginStepEvent)
	}
}

func (e *Engine) applyBeginGameEvent(
	game state.Game,
	beginGameEvent event.BeginGameEvent,
) (state.Game, error) {
	e.log.Info("Game started")
	return game, nil
}

func (e *Engine) applyEndGameEvent(
	game state.Game,
	evnt event.EndGameEvent,
) (state.Game, error) {
	e.log.Info("Game over, winner:", evnt.WinnerID)
	game = game.WithGameOver(evnt.WinnerID)
	return game, mtg.ErrGameOver
}

func (e *Engine) applyEndStepEvent(
	game state.Game,
	endStepEvent event.EndStepEvent,
) (state.Game, error) {
	game = game.WithClearedPriority().WithResetPriorityPasses()
	for _, playerID := range game.PlayerIDsInTurnOrder() {
		player, ok := game.GetPlayer(playerID)
		if !ok {
			return game, fmt.Errorf("player %q not found", playerID)
		}
		player = player.WithEmptyManaPool()

		game = game.WithUpdatedPlayer(player)
	}
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
	return game, fmt.Errorf("unknown end step event '%T'", endStepEvent)
}

func (e *Engine) applyBeginTurnEvent(game state.Game, evnt event.BeginTurnEvent) (state.Game, error) {
	// Reset player state for the new turn
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	player = player.WithNextTurn()
	game = game.WithUpdatedPlayer(player)
	battlefield := game.Battlefield().RemoveSummoningSickness(evnt.PlayerID)
	game = game.WithBattlefield(battlefield)
	return game, nil
}

func (e *Engine) applyEndTurnEvent(game state.Game, evnt event.EndTurnEvent) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	player = player.WithClearSpellsCastsThisTurn()
	player = player.WithClearLandPlayedThisTurn()
	game = game.WithUpdatedPlayer(player)
	return game, nil
}
