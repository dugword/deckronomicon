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

// TODO: There's probably a generic "PutBackInZone" effect that this could be a special case of that.

func PutBackOnTopHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []definition.EffectModifier,
) (EffectResult, error) {
	putBackCount := 1
	for _, modifier := range modifiers {
		if modifier.Key == "Count" {
			count, err := strconv.Atoi(modifier.Value)
			if err != nil {
				return EffectResult{}, fmt.Errorf("invalid modifier %q for PutBackOnTop effect: %w", modifier.Key, err)
			}
			putBackCount = count
		}
	}
	if putBackCount == 0 {
		return EffectResult{}, fmt.Errorf("invalid required modifier %q for PutBackOnTop effect", "Count")
	}
	choicePrompt := choose.ChoicePrompt{
		Message: fmt.Sprintf("Put %d card(s) on top of your library in any order", putBackCount),
		Source:  source,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(player.Hand().GetAll()),
			Min:     putBackCount,
			Max:     putBackCount,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for PutBackOnTop")
		}
		if len(selected.Choices) != putBackCount {
			return EffectResult{}, errors.New("incorrect number of choices selected for PutBackOnTop")
		}
		var events []event.GameEvent
		for _, choice := range selected.Choices {
			card, ok := choice.(gob.Card)
			if !ok {
				return EffectResult{}, errors.New("selected choice is not a card in a zone")
			}
			events = append(events, event.PutCardOnTopOfLibraryEvent{
				PlayerID: player.ID(),
				CardID:   card.ID(),
				FromZone: mtg.ZoneHand,
			})
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
