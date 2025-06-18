package engine

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"fmt"
)

func (e *Engine) ResolveEffect(
	game state.Game,
	player state.Player,
	resolvable state.Resolvable,
	target target.TargetValue,
	effectSpec definition.EffectSpec,
) ([]event.GameEvent, error) {
	var events []event.GameEvent
	efct, err := effect.Build(effectSpec)
	if err != nil {
		return nil, fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
	}
	effectResult, err := efct.Resolve(e.game, player, resolvable, target)
	if err != nil {
		return nil, fmt.Errorf("failed to apply effect %q: %w", effectSpec.Name, err)
	}
	// apply events as we go, instead of at the end.
	// later effects might depend on earlier ones.
	// TODO: This whole section feels clunky.
	for _, evnt := range effectResult.Events {
		if err := e.ApplyEvent(evnt); err != nil {
			return nil, fmt.Errorf("failed to apply event %T: %w", evnt, err)
		}
	}
	// TODO: This needs to be a loop, right now we only handle a depth of 1
	if effectResult.ChoicePrompt.ChoiceOpts != nil {
		agent := e.agents[player.ID()]
		choiceResults, err := agent.Choose(effectResult.ChoicePrompt)
		if err != nil {
			return nil, fmt.Errorf("failed to choose action for player %q: %w", player.ID(), err)
		}
		if effectResult.ResumeFunc == nil {
			return nil, fmt.Errorf("effect %q requires choices but has no resume function", efct.Name())
		}
		effectResult, err = effectResult.ResumeFunc(choiceResults)
		if err != nil {
			return nil, fmt.Errorf("failed to resume function for effect %s: %w", efct.Name(), err)
		}
		if effectResult.ChoicePrompt.ChoiceOpts != nil {
			panic("only one level of recursion supported")
		}
		events = append(events, effectResult.Events...)
	}
	return events, nil
}
