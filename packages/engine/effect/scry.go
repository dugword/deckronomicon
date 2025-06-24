package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/target"

	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type ScryEffect struct {
	Count int `json:"Count"`
}

func NewScryEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var scryEffect ScryEffect
	count, ok := effectSpec.Modifiers["Count"].(int)
	if !ok || count <= 0 {
		return nil, fmt.Errorf("ScryEffect requires a 'Count' modifier of type int greater than 0, got %T", effectSpec.Modifiers["Count"])
	}
	scryEffect.Count = count
	return scryEffect, nil
}

func (e ScryEffect) Name() string {
	return "Scry"
}

func (e ScryEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e ScryEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	if e.Count <= 0 {
		return EffectResult{}, fmt.Errorf("invalid required modifier %q for Scry effect", "Count")
	}
	cards, _ := player.Library().TakeN(e.Count)
	choicePrompt := choose.ChoicePrompt{
		Message: "Put each card on top or bottom of your library in any order",
		Source:  source,
		ChoiceOpts: choose.MapChoicesToBucketsOpts{
			Choices: choose.NewChoices(cards),
			Buckets: []choose.Bucket{
				choose.BucketTop,
				choose.BucketBottom,
			},
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.MapChoicesToBucketsResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for Scrying")
		}
		if len(selected.Assignments) == 0 {
			return EffectResult{}, errors.New("no choices selected for Scrying")
		}
		var events []event.GameEvent
		for bucket, choices := range selected.Assignments {
			if bucket != choose.BucketTop && bucket != choose.BucketBottom {
				return EffectResult{}, fmt.Errorf("invalid bucket %q for Scrying", bucket)
			}
			for _, choice := range choices {
				card, ok := choice.(gob.Card)
				if !ok {
					return EffectResult{}, errors.New("choice is not a card")
				}
				switch bucket {
				case choose.BucketTop:
					events = append(events, event.PutCardOnTopOfLibraryEvent{
						PlayerID: player.ID(),
						CardID:   card.ID(),
						FromZone: mtg.ZoneLibrary,
					})
				case choose.BucketBottom:
					events = append(events, event.PutCardOnBottomOfLibraryEvent{
						PlayerID: player.ID(),
						CardID:   card.ID(),
						FromZone: mtg.ZoneLibrary,
					})
				default:
					return EffectResult{}, fmt.Errorf("invalid bucket %q for Scrying", bucket)
				}
			}
		}
		return EffectResult{
			Events: events,
		}, nil
	}
	return EffectResult{
		ChoicePrompt: choicePrompt,
		ResumeFunc:   resumeFunc,
	}, nil
}
