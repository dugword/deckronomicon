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

func ResolveScry(
	game *state.Game,
	playerID string,
	scry *effect.Scry,
	source gob.Object,
) (Result, error) {
	player := game.GetPlayer(playerID)
	cards := player.Library().GetN(scry.Count)
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
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		selected, ok := choiceResults.(choose.MapChoicesToBucketsResults)
		if !ok {
			return Result{}, errors.New("invalid choice results for Scrying")
		}
		if len(selected.Assignments) == 0 {
			return Result{}, errors.New("no choices selected for Scrying")
		}
		var events []event.GameEvent
		for bucket, choices := range selected.Assignments {
			if bucket != choose.BucketTop && bucket != choose.BucketBottom {
				return Result{}, fmt.Errorf("invalid bucket %q for Scrying", bucket)
			}
			for _, choice := range choices {
				card, ok := choice.(*gob.Card)
				if !ok {
					return Result{}, errors.New("choice is not a card")
				}
				switch bucket {
				case choose.BucketTop:
					events = append(events, &event.PutCardOnTopOfLibraryEvent{
						PlayerID: player.ID(),
						CardID:   card.ID(),
						FromZone: mtg.ZoneLibrary,
					})
				case choose.BucketBottom:
					events = append(events, &event.PutCardOnBottomOfLibraryEvent{
						PlayerID: player.ID(),
						CardID:   card.ID(),
						FromZone: mtg.ZoneLibrary,
					})
				default:
					return Result{}, fmt.Errorf("invalid bucket %q for Scrying", bucket)
				}
			}
		}
		return Result{
			Events: events,
		}, nil
	}
	return Result{
		ChoicePrompt: choicePrompt,
		Resume:       resumeFunc,
	}, nil
}
