package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
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
) (EffectResult, error) {
	var events []event.GameEvent
	for range e.Count {
		events = append(events, event.DiscardCardEvent{
			PlayerID: player.ID(),
		})
	}
	return EffectResult{
		Events: events,
	}, nil
}
