package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type AdditionalManaEffect struct {
	Subtype  mtg.Subtype  `json:"Subtype"`
	Mana     string       `json:"Mana"`
	Duration mtg.Duration `json:"Duration"`
}

func NewAdditionalManaEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var additionalManaEffect AdditionalManaEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &additionalManaEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AdditionalManaEffect: %w", err)
	}
	return additionalManaEffect, nil
}

func (e AdditionalManaEffect) Name() string {
	return "AdditionalMana"
}

func (e AdditionalManaEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e AdditionalManaEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	evnt := event.RegisterTriggeredEffectEvent{
		PlayerID: player.ID(),
		Trigger: state.Trigger{
			EventType: "LandTappedForMana",
			Filter: state.Filter{
				Subtypes: []mtg.Subtype{e.Subtype},
			},
		},
		Duration: e.Duration,
		EffectSpecs: []definition.EffectSpec{
			{
				Name:      "AddMana",
				Modifiers: json.RawMessage(fmt.Sprintf(`{"Mana": "%s"}`, e.Mana)),
			},
		},
	}
	return EffectResult{
		Events: []event.GameEvent{evnt},
	}, nil
}
