package resolver

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func ResolveShuffleFromGraveyard(
	game state.Game,
	playerID string,
	shuffleFromGraveyard effect.ShuffleFromGraveyard,
	source gob.Object,
	resEnv *resenv.ResEnv,
) (Result, error) {
	player := game.GetPlayer(playerID)
	cards := player.Graveyard().GetAll()
	if len(cards) == 0 {
		var cardIDs []string
		for _, card := range player.Library().GetAll() {
			cardIDs = append(cardIDs, card.ID())
		}
		shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
		return Result{
			Events: []event.GameEvent{
				event.ShuffleLibraryEvent{
					PlayerID:         player.ID(),
					ShuffledCardsIDs: shuffledCardsIDs,
				},
			},
		}, nil
	}
	choicePrompt := choose.ChoicePrompt{
		Message: "Choose cards to shuffle into your library",
		Source:  source,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(cards),
			Max:     shuffleFromGraveyard.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		var events []event.GameEvent
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return Result{}, errors.New("invalid choice results for shuffling from graveyard")
		}
		var cardIDs []string
		for _, card := range player.Library().GetAll() {
			cardIDs = append(cardIDs, card.ID())
		}
		for _, choice := range selected.Choices {
			card, ok := player.Graveyard().Get(choice.ID())
			if !ok {
				return Result{}, fmt.Errorf("card %q not found in graveyard", choice.ID())
			}
			events = append(events, event.PutCardOnBottomOfLibraryEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneGraveyard,
			})
			// TODO: This doesn't feel super clear,
			// but we need the card ID order to shuffle the library.
			// I think I'd rather have the the cards added to the library
			// and then shuffled.
			// Maybe that could be a separate effect?
			// If we do that I'd need to update how the effectResult chanins
			// are handled. Like just a bool "should continue" or something.
			// This could be useful for other effects where I need to
			// apply state changes before the next part of the effect
			// is resolved.
			cardIDs = append(cardIDs, card.ID())
		}
		shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
		events = append(events, event.ShuffleLibraryEvent{
			PlayerID:         player.ID(),
			ShuffledCardsIDs: shuffledCardsIDs,
		})
		return Result{
			Events: events,
		}, nil
	}
	return Result{
		ChoicePrompt: choicePrompt,
		Resume:       resumeFunc,
	}, nil
}
