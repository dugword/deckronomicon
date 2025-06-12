package gob

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"fmt"
)

type AbilityOnStack struct {
	id          string
	name        string
	owner       string
	constroller string
	sourceID    string
	abilityID   string
	effects     []Effect
}

func (a AbilityOnStack) Description() string {
	return fmt.Sprintf("Write a better description: %s", a.abilityID)
}

func (a AbilityOnStack) Effects() []Effect {
	return a.effects
}

func (a AbilityOnStack) Name() string {
	return a.name
}

func (a AbilityOnStack) Match(predicate query.Predicate) bool {
	return predicate(a)
}

func (a AbilityOnStack) ID() string {
	return a.id
}

func (a AbilityOnStack) Owner() string {

	return a.owner
}
func (a AbilityOnStack) Controller() string {
	return a.constroller
}

func NewAbilityOnStack(id string,
	playerID string,
	sourceID string,
	abilityID string,
	abilityName string,
	effectSpecs []definition.EffectSpec,
) AbilityOnStack {
	abilityOnStack := AbilityOnStack{
		id:          id,
		owner:       playerID,
		constroller: playerID,
		sourceID:    sourceID,
		abilityID:   abilityID,
		name:        abilityName,
	}
	var effects []Effect
	// This is redundant with the code in load_from_definitions, it should probably be removed from there
	for _, effectSpec := range effectSpecs {
		var modifiers []Tag
		for _, modifier := range effectSpec.Modifiers {
			// TODO: Check if the modifier is valid for the effect
			modifier := Tag{
				Key:   modifier.Key,
				Value: modifier.Value,
			}
			modifiers = append(modifiers, modifier)
		}
		// TODO: Check if required modifiers are present
		var tags []Tag
		// TODO load the tags from the effect spec
		effect := Effect{
			// TODO: Check if the effect exists
			name:      effectSpec.Name,
			modifiers: modifiers,
			optional:  effectSpec.Optional,
			tags:      tags,
		}
		effects = append(effects, effect)
	}
	abilityOnStack.effects = effects
	return abilityOnStack
}
