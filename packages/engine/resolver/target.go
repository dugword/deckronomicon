package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

func ResolveTarget(
	game *state.Game,
	playerID string,
) (Result, error) {
	return Result{
		Events: []event.GameEvent{},
	}, nil
}
