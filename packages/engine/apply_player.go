package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// Player Events are events that are triggered by player actions. These are mostly to log the player's
// intent, and the actual state changes are handled in the state change events. However, some player events
// do set flags (like land played for turn) or increment internal counters (like storm count). Or mark
// creates as attacking/blocking, etc.
// Changes that represent "visable" state changes should be handled in the state change events.

func (e *Engine) applyPlayerEvent(game state.Game, playerEvent event.PlayerEvent) (state.Game, error) {
	switch evnt := playerEvent.(type) {
	case event.ActivateAbilityEvent:
		return game, nil
	case event.AssignCombatDamageEvent:
		return game, nil
	case event.CastSpellEvent:
		return e.applyCastSpellEvent(game, evnt)
	case event.ConcedeEvent:
		return game, nil
	case event.CycleCardEvent:
		return game, nil
	case event.DeclareAttackersEvent:
		return game, nil
	case event.DeclareBlockersEvent:
		return game, nil
	case event.PassPriorityEvent:
		return e.applyPassPriorityEvent(game, evnt)
	case event.PlayLandEvent:
		return e.applyPlayLandEvent(game, evnt)
	case event.ClearRevealedEvent:
		return e.applyClearRevealedEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown player event type '%T'", evnt)
	}
}

func (e *Engine) applyCastSpellEvent(
	game state.Game,
	event event.CastSpellEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(event.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", event.PlayerID)
	}
	player = player.WithSpellCastThisTurn()
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyPassPriorityEvent(
	game state.Game,
	evnt event.PassPriorityEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	game = game.WithPlayerPassedPriority(
		evnt.PlayerID,
	)
	game = game.WithUpdatedPlayer(player)
	return game, nil
}

func (e *Engine) applyPlayLandEvent(
	game state.Game,
	evnt event.PlayLandEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	game = game.WithUpdatedPlayer(player.WithLandPlayedThisTurn())
	return game, nil
}

func (e *Engine) applyClearRevealedEvent(
	game state.Game,
	evnt event.ClearRevealedEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	game = game.WithUpdatedPlayer(player.WithClearRevealed())
	return game, nil
}
