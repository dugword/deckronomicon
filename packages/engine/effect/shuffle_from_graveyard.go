package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"errors"
	"fmt"
)

type ShuffleFromGraveyardEffect struct {
	Count int `json:"Count,omitempty"`
}

func NewShuffleFromGraveyardEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var shuffleFromGraveyardEffect ShuffleFromGraveyardEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &shuffleFromGraveyardEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ShuffleFromGraveyardEffect: %w", err)
	}
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
) (EffectResult, error) {
	cards := player.Graveyard().GetAll()
	if len(cards) == 0 {
		return EffectResult{
			Events: []event.GameEvent{
				event.ShuffleLibraryEvent{
					PlayerID: player.ID(),
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
		}
		events = append(events, event.ShuffleLibraryEvent{
			PlayerID: player.ID(),
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
