package effect

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"fmt"
)

type RegisterDelayedTriggeredAbility struct {
	EventType string
	Filter    query.Opts
	Effects   []Effect
}

func (e *RegisterDelayedTriggeredAbility) Name() string {
	return "RegisterDelayedTriggeredAbility"
}

func (e *RegisterDelayedTriggeredAbility) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func NewRegisterDelayedTriggeredAbility(modifier map[string]any) (*RegisterDelayedTriggeredAbility, error) {
	eventType, ok := modifier["EventType"].(string)
	if !ok {
		return nil, fmt.Errorf("a 'EventType' modifier of type string required, got %T", modifier["EventType"])
	}
	eventFilterRaw, ok := modifier["EventFilter"].(map[string]any)
	if eventFilterRaw != nil && !ok {
		return nil, fmt.Errorf("a 'EventFilter' modifier of type map[string]any required, got %T", modifier["EventFilter"])
	}
	query, err := buildQueryOpts(eventFilterRaw)
	if err != nil {
		return nil, err
	}
	effectsRaw, ok := modifier["Effects"].([]any)
	if !ok {
		return nil, fmt.Errorf("a Effects' modifier of type []any required, got %T", modifier["Effects"])
	}
	effects, err := parseEffects(effectsRaw)
	return &RegisterDelayedTriggeredAbility{
		EventType: eventType,
		Filter:    query,
		Effects:   effects,
	}, nil
}

func parseEffects(raw []any) ([]Effect, error) {
	var effects []Effect
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
			effectSpec := definition.Effect{
				Name:      name,
				Modifiers: modifiers,
			}
			effect, err := New(&effectSpec)
			if err != nil {
				return nil, fmt.Errorf("error creating effect %s: %w", name, err)
			}
			effects = append(effects, effect)
		}
	}
	return effects, nil
}
