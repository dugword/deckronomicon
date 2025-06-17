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

// TODO: Make the names consistent, registered delayed effect vs delayed trigger vs ability or I dunno.
type RegisterDelayedEffectEffect struct {
	Trigger state.Trigger           `json:"Trigger"`
	Effects []definition.EffectSpec `json:"Effects"`
}

func NewRegisterDelayedEffectEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var RegisterDelayedEffectEffect RegisterDelayedEffectEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &RegisterDelayedEffectEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RegisterDelayedEffectEffect: %w", err)
	}
	return RegisterDelayedEffectEffect, nil
}

func (e RegisterDelayedEffectEffect) Name() string {
	return "RegisterDelayedEffect"
}

func (e RegisterDelayedEffectEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e RegisterDelayedEffectEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	events := []event.GameEvent{
		// TODO: add source to the event, or source name.
		event.RegisterTriggeredEffectEvent{
			PlayerID:   player.ID(),
			SourceName: source.Name(),
			SourceID:   source.ID(),
			Trigger: state.Trigger{
				EventType: e.Trigger.EventType,
			},
			OneShot:     true,
			EffectSpecs: e.Effects,
		},
	}
	return EffectResult{
		Events: events,
	}, nil
}
