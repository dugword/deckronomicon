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

func parseFindCardCheatCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (action.FindCardCheatAction, error) {
	player := game.GetPlayer(playerID)
	cards := player.Library().GetAll()
	var card *gob.Card
	var err error
	if idOrName == "" {
		card, err = buildFindCardCheatCommandByChoice(cards, playerID, agent)
		if err != nil {
			return action.FindCardCheatAction{}, fmt.Errorf("failed to choose a card to find: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.FindCardCheatAction{}, fmt.Errorf("no card found with ID or name %q", idOrName)
		}
		card = found
	}
	return action.NewFindCardCheatAction(card.ID()), nil
}

func buildFindCardCheatCommandByChoice(
	cards []*gob.Card,
	playerID string,
	agent engine.PlayerAgent,
) (*gob.Card, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a card to put into your hand",
		Source:   nil, // TODO: Make this better
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	choiceResults, err := agent.Choose(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected a single choice result")
	}
	card, ok := selected.Choice.(*gob.Card)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card")
	}
	return card, nil
}
