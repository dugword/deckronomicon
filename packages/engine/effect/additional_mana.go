package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

type AdditionalManaEffect struct {
	Subtype  mtg.Subtype
	Mana     string
	Duration mtg.Duration
}

func NewAdditionalManaEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var additionalManaEffect AdditionalManaEffect
	subtypeString, ok := effectSpec.Modifiers["Subtype"].(string)
	if !ok {
		return nil, fmt.Errorf("AdditionalManaEffect requires a 'Subtype' modifier of type string, got %T", effectSpec.Modifiers["Subtype"])
	}
	subtype, ok := mtg.StringToSubtype(subtypeString)
	if !ok {
		return nil, fmt.Errorf("AdditionalManaEffect requires a valid 'Subtype' modifier, got %q", subtypeString)
	}
	if subtype == "" {
		return nil, fmt.Errorf("AdditionalManaEffect requires a non-empty 'Subtype' modifier")
	}
	additionalManaEffect.Subtype = subtype
	manaString, ok := effectSpec.Modifiers["Mana"].(string)
	if !ok {
		return nil, fmt.Errorf("AdditionalManaEffect requires a 'Mana' modifier of type string, got %T", effectSpec.Modifiers["Mana"])
	}
	if manaString == "" {
		return nil, fmt.Errorf("AdditionalManaEffect requires a non-empty 'Mana' modifier")
	}
	additionalManaEffect.Mana = manaString
	durationString, ok := effectSpec.Modifiers["Duration"].(string)
	if !ok {
		return nil, fmt.Errorf("AdditionalManaEffect requires a 'Duration' modifier of type mtg.Duration, got %T", effectSpec.Modifiers["Duration"])
	}
	duration, ok := mtg.StringToDuration(durationString)
	if !ok {
		return nil, fmt.Errorf("AdditionalManaEffect requires a valid 'Duration' modifier, got %q", durationString)
	}
	if duration == "" {
		return nil, fmt.Errorf("AdditionalManaEffect requires a non-empty 'Duration' modifier")
	}
	additionalManaEffect.Duration = duration
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
	resEnv *resenv.ResEnv,
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
				Modifiers: map[string]any{"Mana": e.Mana},
			},
		},
	}
	return EffectResult{
		Events: []event.GameEvent{evnt},
	}, nil
}
