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

// TODO: Make the names consistent, registered delayed effect vs delayed trigger vs ability or I dunno.
type RegisterDelayedEffectEffect struct {
	Trigger state.Trigger
	Effects []definition.EffectSpec
}

func parseEffects(raw []any) ([]definition.EffectSpec, error) {
	var effectSpecs []definition.EffectSpec
	for _, effectRaw := range raw {
		effectSpecRaw, ok := effectRaw.(map[string]any)
		if ok {
			name, ok := effectSpecRaw["Name"].(string)
			if !ok {
				return nil, fmt.Errorf("RegisterDelayedEffectEffect requires each effect in 'Effects' modifier to have a 'Name' key of type string, got %T", effectSpecRaw["Name"])
			}
			modifiers, ok := effectSpecRaw["Modifiers"].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("RegisterDelayedEffectEffect requires each effect in 'Effects' modifier to have a 'Modifiers' key of type map[string]any, got %T", effectSpecRaw["Modifiers"])
			}
			effectSpec := definition.EffectSpec{
				Name:      name,
				Modifiers: modifiers,
			}
			effectSpecs = append(effectSpecs, effectSpec)
		}
	}
	return effectSpecs, nil
}

func parseTrigger(raw map[string]any) (state.Trigger, error) {
	trigger := state.Trigger{}
	eventType, ok := raw["EventType"].(string)
	if !ok {
		return trigger, fmt.Errorf("RegisterDelayedEffectEffect requires 'Trigger' modifier to have an 'EventType' key of type string, got %T", raw["EventType"])
	}
	trigger.EventType = eventType
	filterRaw, ok := raw["Filter"].(map[string]any)
	if !ok {
		filterRaw = map[string]any{}
	}
	cardTypesRaw, ok := filterRaw["CardTypes"].([]any)
	if ok {
		for _, cardTypeRaw := range cardTypesRaw {
			cardType, ok := mtg.StringToCardType(fmt.Sprintf("%v", cardTypeRaw))
			if !ok {
				return trigger, fmt.Errorf("RegisterDelayedEffectEffect requires 'Filter' modifier to have a 'CardTypes' key of type []string, got %T", cardTypeRaw)
			}
			trigger.Filter.CardTypes = append(trigger.Filter.CardTypes, cardType)
		}
	}
	subtypesRaw, ok := filterRaw["Subtypes"].([]any)
	if ok {
		for _, subtypeRaw := range subtypesRaw {
			subtype, ok := mtg.StringToSubtype(fmt.Sprintf("%v", subtypeRaw))
			if !ok {
				return trigger, fmt.Errorf("RegisterDelayedEffectEffect requires 'Filter' modifier to have a 'Subtypes' key of type []string, got %T", subtypeRaw)
			}
			trigger.Filter.Subtypes = append(trigger.Filter.Subtypes, subtype)
		}
	}
	return trigger, nil
}

// TODO: This is SUPER clunky - make it less bad
func NewRegisterDelayedEffectEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var registerDelayedEffectEffect RegisterDelayedEffectEffect
	triggerRaw, ok := effectSpec.Modifiers["Trigger"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("RegisterDelayedEffectEffect requires 'Trigger' modifier to be a map, got %T", effectSpec.Modifiers["Trigger"])
	}
	trigger, err := parseTrigger(triggerRaw)
	if err != nil {
		return nil, err
	}
	registerDelayedEffectEffect.Trigger = trigger
	effectsRaw, ok := effectSpec.Modifiers["Effects"].([]any)
	if !ok {
		return nil, fmt.Errorf("RegisterDelayedEffectEffect requires 'Effects' modifier to be a slice, got %T", effectSpec.Modifiers["Effects"])
	}
	effects, err := parseEffects(effectsRaw)
	if err != nil {
		return nil, err
	}
	registerDelayedEffectEffect.Effects = effects
	return registerDelayedEffectEffect, nil
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
	resEnv *resenv.ResEnv,
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
