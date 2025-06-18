package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

func parseViewCommand(
	arg string,
	game state.Game,
	player state.Player,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (action.ViewAction, error) {
	return action.NewViewAction(player, "", arg), nil
}
