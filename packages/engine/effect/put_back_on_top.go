package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/target"

	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"errors"
	"fmt"
)

type PutBackOnTopEffect struct {
	Count int `json:"Count"`
}

func NewPutBackOnTopEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var putBackOnTopEffect PutBackOnTopEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &putBackOnTopEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PutBackOnTopEffectModifiers: %w", err)
	}
	return putBackOnTopEffect, nil
}

func (e PutBackOnTopEffect) Name() string {
	return "PutBackOnTop"
}

func (e PutBackOnTopEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

// TODO: There's probably a generic "PutBackInZone" effect that this could be a special case of that.

func (e PutBackOnTopEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	targets target.TargetValue,
) (EffectResult, error) {
	if e.Count == 0 {
		return EffectResult{}, fmt.Errorf("invalid required modifier %q for PutBackOnTop effect", "Count")
	}
	choicePrompt := choose.ChoicePrompt{
		Message: fmt.Sprintf("Put %d card(s) on top of your library in any order", e.Count),
		Source:  source,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(player.Hand().GetAll()),
			Min:     e.Count,
			Max:     e.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for PutBackOnTop")
		}
		if len(selected.Choices) != e.Count {
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
