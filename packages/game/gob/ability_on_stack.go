package gob

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/query"
	"fmt"
)

type AbilityOnStack struct {
	abilityID         string
	controller        string
	effectWithTargets []effect.EffectWithTarget
	id                string
	name              string
	owner             string
	sourceID          string
}

func (a AbilityOnStack) Controller() string {
	return a.controller
}

func (a AbilityOnStack) Description() string {
	return fmt.Sprintf("Write a better description: %s", a.abilityID)
}

func (a AbilityOnStack) EffectWithTargets() []effect.EffectWithTarget {
	return a.effectWithTargets
}

func (a AbilityOnStack) ID() string {
	return a.id
}

func (a AbilityOnStack) Name() string {
	return a.name
}

func (a AbilityOnStack) Match(predicate query.Predicate) bool {
	return predicate(a)
}

func (a AbilityOnStack) Owner() string {

	return a.owner
}

func (a AbilityOnStack) SourceID() string {
	return a.sourceID
}

func NewAbilityOnStack(id string,
	playerID string,
	sourceID string,
	abilityID string,
	abilityName string,
	effectWithTargets []effect.EffectWithTarget,
) AbilityOnStack {
	abilityOnStack := AbilityOnStack{
		id:                id,
		owner:             playerID,
		controller:        playerID,
		sourceID:          sourceID,
		abilityID:         abilityID,
		name:              abilityName,
		effectWithTargets: effectWithTargets,
	}
	return abilityOnStack
}
