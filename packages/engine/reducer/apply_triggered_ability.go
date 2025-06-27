package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func applyTriggeredAbilityEvent(game state.Game, triggeredEvent event.TriggeredAbilityEvent) (state.Game, error) {
	switch evnt := triggeredEvent.(type) {
	case event.RegisterTriggeredAbilityEvent:
		return applyRegisterTriggeredAbilityEvent(game, evnt)
	case event.RemoveTriggeredAbilityEvent:
		return applyRemoveTriggeredAbilityEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown triggered event type '%T'", evnt)
	}
}

func applyRegisterTriggeredAbilityEvent(
	game state.Game,
	evnt event.RegisterTriggeredAbilityEvent,
) (state.Game, error) {
	player := game.GetPlayer(evnt.PlayerID)
	game = game.WithRegisteredTriggeredAbility(
		player.ID(),
		evnt.SourceName,
		evnt.SourceID,
		evnt.Trigger,
		evnt.Effects,
		evnt.Duration,
		evnt.OneShot,
	)
	return game, nil
}

func applyRemoveTriggeredAbilityEvent(
	game state.Game,
	evnt event.RemoveTriggeredAbilityEvent,
) (state.Game, error) {
	game = game.WithRemoveTriggeredAbility(evnt.ID)
	return game, nil
}
