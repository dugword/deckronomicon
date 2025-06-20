package effect

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"errors"
	"fmt"
)

type DiscardEffect struct {
	Count  int    `json:"Count"`
	Target string `json:"Target"`
}

func NewDiscardEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var discardEffect DiscardEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &discardEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DiscardEffectModifiers: %w", err)
	}
	return discardEffect, nil
}

func (d DiscardEffect) Name() string {
	return "Discard"
}

func (d DiscardEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		return target.NoneTargetSpec{}
	case "Player":
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DiscardEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e DiscardEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	if e.Count <= 0 {
		return EffectResult{}, fmt.Errorf("invalid required modifier %q for Discard effect", "Count")
	}
	cards := player.Hand().GetAll()
	choicePrompt := choose.ChoicePrompt{
		Message: "Chose cards to discard",
		Source:  source,
		ChoiceOpts: choose.ChooseManyOpts{
			Choices: choose.NewChoices(cards),
			Min:     e.Count,
			Max:     e.Count,
		},
	}
	resumeFunc := func(choiceResults choose.ChoiceResults) (EffectResult, error) {
		selected, ok := choiceResults.(choose.ChooseManyResults)
		if !ok {
			return EffectResult{}, errors.New("invalid choice results for Discarding")
		}
		var events []event.GameEvent
		for _, choice := range selected.Choices {
			// Create the discard event
			discardEvent := event.DiscardCardEvent{
				PlayerID: player.ID(),
				CardID:   choice.ID(),
			}
			// Apply the discard event to the game state
			events = append(events, discardEvent)
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
