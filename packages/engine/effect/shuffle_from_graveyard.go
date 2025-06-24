package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type ShuffleFromGraveyardEffect struct {
	Count int
}

func NewShuffleFromGraveyardEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var shuffleFromGraveyardEffect ShuffleFromGraveyardEffect
	count, ok := effectSpec.Modifiers["Count"].(int)
	if !ok || count <= 0 {
		return nil, fmt.Errorf("ShuffleFromGraveyardEffect requires a 'Count' modifier of type int greater than 0, got %T", effectSpec.Modifiers["Count"])
	}
	shuffleFromGraveyardEffect.Count = count
	return shuffleFromGraveyardEffect, nil
}

func (e ShuffleFromGraveyardEffect) Name() string {
	return "ShuffleFromGraveyard"
}

func (e ShuffleFromGraveyardEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e ShuffleFromGraveyardEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	cards := player.Graveyard().GetAll()
	if len(cards) == 0 {
		var cardIDs []string
		for _, card := range player.Library().GetAll() {
			cardIDs = append(cardIDs, card.ID())
		}
		shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
		return EffectResult{
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
			Max:     e.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		var events []event.GameEvent
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for shuffling from graveyard")
		}
		var cardIDs []string
		for _, card := range player.Library().GetAll() {
			cardIDs = append(cardIDs, card.ID())
		}
		for _, choice := range selected.Choices {
			card, ok := player.Graveyard().Get(choice.ID())
			if !ok {
				return EffectResult{}, fmt.Errorf("card %q not found in graveyard", choice.ID())
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
		return EffectResult{
			Events: events,
		}, nil
	}
	return EffectResult{
		ChoicePrompt: choicePrompt,
		ResumeFunc:   resumeFunc,
	}, nil
}
