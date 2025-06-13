package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"strconv"
)

func ScryEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	scryCount := 1
	for _, modifier := range modifiers {
		if modifier.Key == "Count" {
			count, err := strconv.Atoi(modifier.Value)
			if err != nil {
				return EffectResult{}, fmt.Errorf("invalid modifier %q for Scry effect: %w", modifier.Key, err)
			}
			scryCount = count
		}
	}
	if scryCount == 0 {
		return EffectResult{}, fmt.Errorf("missing required modifier %q for Scry effect", "Count")
	}
	cards, _ := player.Library().TakeN(scryCount)
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
