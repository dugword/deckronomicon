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
	case event.PassPriorityEvent:
		return applyPassPriorityEvent(game, evnt)
	case event.ReceivePriorityEvent:
		return game, nil
	case event.ResetPriorityPassesEvent:
		return applyResetPriorityPassesEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown priority event type '%T'", evnt)
	}
}

func applyPassPriorityEvent(
	game state.Game,
	evnt event.PassPriorityEvent,
) (state.Game, error) {
	game = game.WithPlayerPassedPriority(evnt.PlayerID)
	return game, nil
}

func applyResetPriorityPassesEvent(
	game state.Game,
	_ event.ResetPriorityPassesEvent,
) (state.Game, error) {
	game = game.WithResetPriorityPasses()
	return game, nil
}
