package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseDiscardCheatCommand(
	idOrName string,
	game state.Game,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (action.DiscardCheatAction, error) {
	cards := player.Hand().GetAll()
	var card gob.Card
	var err error
	if idOrName == "" {
		card, err = buildDiscardCommandByChoice(cards, player, chooseFunc)
		if err != nil {
			return action.DiscardCheatAction{}, fmt.Errorf("failed to choose a card to discard: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.DiscardCheatAction{}, fmt.Errorf("no land found with ID or name %q", idOrName)
		}
		card = found
	}
	return action.NewDiscardCheatAction(player, card), nil
}

func buildDiscardCommandByChoice(
	cards []gob.Card,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (gob.Card, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a card to discard",
		Source:   player,
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return gob.Card{}, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return gob.Card{}, fmt.Errorf("expected a single choice result")
	}
	card, ok := selected.Choice.(gob.Card)
	if !ok {
		return gob.Card{}, fmt.Errorf("selected choice is not a card in hand")
	}
	return card, nil
}
