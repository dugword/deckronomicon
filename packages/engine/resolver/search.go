package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func ResolveSearch(
	game state.Game,
	playerID string,
	search effect.Search,
	source gob.Object,
) (Result, error) {
	player := game.GetPlayer(playerID)
	query, err := buildPredicate(QueryOpts(search))
	if err != nil {
		return Result{}, fmt.Errorf("failed to build query for Search effect: %w", err)
	}
	cards := player.Library().FindAll(query)
	choicePrompt := choose.ChoicePrompt{
		// TODO: provide more detail on what kind of card to choose
		Message: "Choose a card to put into your hand",
		Source:  source,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(cards),
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		selected, ok := choiceResults.(choose.ChooseOneResults)
		if !ok {
			return Result{}, fmt.Errorf("expected a single choice result")
		}
		card, ok := selected.Choice.(gob.Card)
		if !ok {
			return Result{}, errors.New("choice is not a card")
		}
		events := []event.GameEvent{
			event.PutCardInHandEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneLibrary,
			},
		}
		return Result{
			Events: events,
		}, nil
	}
	// Need to get choices
	return Result{
		ChoicePrompt: choicePrompt,
		Resume:       resumeFunc,
	}, nil
}
