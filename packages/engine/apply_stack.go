package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func (e *Engine) applyStackEvent(game state.Game, stackEvent event.StackEvent) (state.Game, error) {
	switch evnt := stackEvent.(type) {
	case event.ResolveTopObjectOnStackEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown stack event type '%T'", evnt)
	}
}
