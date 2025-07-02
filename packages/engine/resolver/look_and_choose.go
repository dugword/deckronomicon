package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/take"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func ResolveLookAndChoose(
	game *state.Game,
	playerID string,
	lookAndChoose *effect.LookAndChoose,
	source gob.Object,
) (Result, error) {
	player := game.GetPlayer(playerID)
	cards, _ := player.Library().TakeN(lookAndChoose.Look)
	var events []event.GameEvent
	for _, card := range cards {
		events = append(events, &event.RevealCardEvent{
			PlayerID: player.ID(),
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
		})
	}
	predicate, err := buildPredicate(QueryOpts{
		CardTypes: lookAndChoose.CardTypes,
	})
	if err != nil {
		return Result{}, fmt.Errorf("failed to build query for LookAndChoose: %w", err)
	}
	choiceCards := query.FindAll(cards, predicate)
	choicePrompt := choose.ChoicePrompt{
		// TODO: Add type information to message
		Message:  fmt.Sprintf("Look at the top %d cards of your library. Choose %d of them", lookAndChoose.Look, lookAndChoose.Choose),
		Source:   source,
		Optional: true,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(choiceCards),
			Max:     lookAndChoose.Choose,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return Result{}, errors.New("invalid choice results for LookAndChoose")
		}
		var selectedCards []*gob.Card
		for _, choice := range selected.Choices {
			taken, remaining, ok := take.By(cards, has.ID(choice.ID()))
			if !ok {
				return Result{}, fmt.Errorf("selected card %q not found in looked at cards", choice.ID())
			}
			selectedCards = append(selectedCards, taken)
			cards = remaining
		}
		var events []event.GameEvent
		for _, card := range selectedCards {
			events = append(events, &event.PutCardInHandEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneLibrary,
			})
		}
		if lookAndChoose.Rest == mtg.ZoneLibrary {
			// Put the rest on the bottom of the library in a random order
			// TODO: Support ordering
			for _, card := range cards {
				events = append(events, &event.PutCardOnBottomOfLibraryEvent{
					PlayerID: player.ID(),
					CardID:   card.ID(),
					FromZone: mtg.ZoneLibrary,
				})
			}
		} else if lookAndChoose.Rest == mtg.ZoneGraveyard {
			// Put the rest in the graveyard
			for _, card := range cards {
				events = append(events, &event.PutCardInGraveyardEvent{
					PlayerID: player.ID(),
					CardID:   card.ID(),
					FromZone: mtg.ZoneLibrary,
				})
			}
		} else {
			return Result{}, fmt.Errorf("invalid Rest zone %q for LookAndChoose", lookAndChoose.Rest)
		}
		return Result{
			Events: events,
		}, nil
	}
	return Result{
		Events:       events,
		ChoicePrompt: choicePrompt,
		Resume:       resumeFunc,
	}, nil
}
