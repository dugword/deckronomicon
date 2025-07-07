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

func ResolvePutBackOnTop(
	game *state.Game,
	playerID string,
	putBackOnTop *effect.PutBackOnTop,
	resolvable state.Resolvable,
) (Result, error) {
	player := game.GetPlayer(playerID)
	choicePrompt := choose.ChoicePrompt{
		Message: fmt.Sprintf("Put %d card(s) on top of your library in any order", putBackOnTop.Count),
		Source:  resolvable,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(player.Hand().GetAll()),
			Min:     putBackOnTop.Count,
			Max:     putBackOnTop.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (Result, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return Result{}, errors.New("invalid choice results for PutBackOnTop")
		}
		if len(selected.Choices) != putBackOnTop.Count {
			return Result{}, errors.New("incorrect number of choices selected for PutBackOnTop")
		}
		var events []event.GameEvent
		for _, choice := range selected.Choices {
			card, ok := choice.(*gob.Card)
			if !ok {
				return Result{}, errors.New("selected choice is not a card in a zone")
			}
			events = append(events, &event.PutCardOnTopOfLibraryEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneHand,
			})
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
