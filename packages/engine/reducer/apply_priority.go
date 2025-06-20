package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

// These are events that manage the priority system in the game.

func applyPriorityEvent(game state.Game, priorityEvent event.PriorityEvent) (state.Game, error) {
	switch evnt := priorityEvent.(type) {
	case event.AllPlayersPassedPriorityEvent:
		return game, nil
	case event.ReceivePriorityEvent:
		newGame := game.WithPlayerWithPriority(
			evnt.PlayerID,
		)
		return newGame, nil
	case event.ResetPriorityPassesEvent:
		newGame := game.WithResetPriorityPasses()
		return newGame, nil
	default:
		return game, fmt.Errorf("unknown priority event type '%T'", evnt)
	}
}
