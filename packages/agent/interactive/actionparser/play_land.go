package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func parsePlayLandCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (action.PlayLandAction, error) {
	ruling := judge.Ruling{Explain: true}
	cards := judge.GetLandsAvailableToPlay(game, playerID, &ruling)
	var cardInZone *gob.CardInZone
	var err error
	if idOrName == "" {
		cardInZone, err = buildPlayLandCommandByChoice(cards, game, playerID, agent)
		if err != nil {
			return action.PlayLandAction{}, fmt.Errorf("failed to choose a land to play: %w", err)
		}
	} else {
		found, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.PlayLandAction{}, fmt.Errorf("no land found with ID or name %q", idOrName)
		}
		cardInZone = found
	}
	request := action.PlayLandRequest{
		CardID: cardInZone.Card().ID(),
	}
	return action.NewPlayLandAction(request), nil
}

func buildPlayLandCommandByChoice(
	cards []*gob.CardInZone,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (*gob.CardInZone, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a land to play",
		Source:   nil,
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
	cardInZone, ok := selected.Choice.(*gob.CardInZone)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in a zone")
	}
	return cardInZone, nil
}
