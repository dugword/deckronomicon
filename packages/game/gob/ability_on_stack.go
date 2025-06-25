package gob

import (
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"fmt"
)

type AbilityOnStack struct {
	id                string
	name              string
	owner             string
	constroller       string
	sourceID          string
	abilityID         string
	effectWithTargets []target.EffectWithTarget
}

func (a AbilityOnStack) Description() string {
	return fmt.Sprintf("Write a better description: %s", a.abilityID)
}

func (a AbilityOnStack) EffectWithTargets() []target.EffectWithTarget {
	return a.effectWithTargets
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

func (a AbilityOnStack) SourceID() string {
	return a.sourceID
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
	effectWithTargets []target.EffectWithTarget,
) AbilityOnStack {
	abilityOnStack := AbilityOnStack{
		id:                id,
		owner:             playerID,
		constroller:       playerID,
		sourceID:          sourceID,
		abilityID:         abilityID,
		name:              abilityName,
		effectWithTargets: effectWithTargets,
	}
	return abilityOnStack
}
