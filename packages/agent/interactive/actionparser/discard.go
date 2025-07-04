package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parseDiscardCheatCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (action.DiscardCheatAction, error) {
	player := game.GetPlayer(playerID)
	cards := player.Hand().GetAll()
	var card *gob.Card
	var err error
	if idOrName == "" {
		card, err = buildDiscardCommandByChoice(cards, game, playerID, agent)
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
	return action.NewDiscardCheatAction(card.ID()), nil
}

func buildDiscardCommandByChoice(
	cards []*gob.Card,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (*gob.Card, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a card to discard",
		Source:   nil, // TODO: Make this better
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := agent.Choose(game, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected a single choice result")
	}
	card, ok := selected.Choice.(*gob.Card)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in hand")
	}
	return card, nil
}
