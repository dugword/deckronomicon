package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

func parseViewCommand(
	arg string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (action.ViewAction, error) {
	return action.NewViewAction("", arg), nil
}
