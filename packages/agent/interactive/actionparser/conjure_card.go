package actionparser

import (
	"deckronomicon/packages/engine/action"
	"errors"
)

func parseConjureCardCheatCommand(
	cardName string,
) (action.ConjureCardCheatAction, error) {
	if cardName == "" {
		return action.ConjureCardCheatAction{}, errors.New("conjure card command requires a card name")
	}
	return action.NewConjureCardCheatAction(cardName), nil
}
