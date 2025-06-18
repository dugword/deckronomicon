package actionparser

import (
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"errors"
)

func parseConjureCardCheatCommand(
	cardName string,
	player state.Player,
) (action.ConjureCardCheatAction, error) {
	if cardName == "" {
		return action.ConjureCardCheatAction{}, errors.New("conjure card command requires a card name")
	}
	return action.NewConjureCardCheatAction(player, cardName), nil
}
