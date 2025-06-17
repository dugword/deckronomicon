package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func applyTriggeredEffect(game state.Game, triggeredEvent event.TriggeredEffectEvent) (state.Game, error) {
	switch evnt := triggeredEvent.(type) {
	case event.RegisterTriggeredEffectEvent:
		return applyRegisterTriggeredEffectEvent(game, evnt)
	case event.RemoveTriggeredEffectEvent:
		return applyRemoveTriggeredEffectEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown triggered event type '%T'", evnt)
	}
}

func applyRegisterTriggeredEffectEvent(
	game state.Game,
	evnt event.RegisterTriggeredEffectEvent,
) (state.Game, error) {
	player, ok := game.GetPlayer(evnt.PlayerID)
	if !ok {
		return game, fmt.Errorf("player %q not found", evnt.PlayerID)
	}
	game = game.WithRegisteredTriggeredEffect(
		player.ID(),
		evnt.SourceName,
		evnt.SourceID,
		evnt.Trigger,
		evnt.EffectSpecs,
		evnt.Duration,
		evnt.OneShot,
	)
	// TODO: Implement the logic for registering the triggered effect
	return game, nil
}

func applyRemoveTriggeredEffectEvent(
	game state.Game,
	evnt event.RemoveTriggeredEffectEvent,
) (state.Game, error) {
	game = game.WithRemoveTriggeredEffect(evnt.ID)
	return game, nil
}
