package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

func ResolveTap(
	game state.Game,
	playerID string,
) (Result, error) {
	var events []event.GameEvent
	return Result{
		Events: events,
	}, nil
}
