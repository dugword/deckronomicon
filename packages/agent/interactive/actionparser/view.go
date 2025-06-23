package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

func parseViewCommand(
	arg string,
	game state.Game,
	player state.Player,
	agent engine.PlayerAgent,
) (action.ViewAction, error) {
	return action.NewViewAction("", arg), nil
}
