package gob

import (
	//	"deckronomicon/packages/game/cost"

	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
	"strings"
)

// TODO Is this used anywhere?

// Ability represents abilities that require activation costs.
type Ability struct {
	name        string
	cost        cost.Cost
	effectSpecs []definition.EffectSpec
	id          string
	zone        mtg.Zone
	source      query.Object
	speed       mtg.Speed
}

func NewAbility(id string) Ability {
	ability := Ability{
		id: id,
	}
	return ability
}

func (a Ability) Controller() string {
	return a.source.Controller()
}

func (a Ability) Owner() string {
	return a.source.Owner()
}

func (a Ability) Cost() cost.Cost {
	return a.cost
}

func (a Ability) EffectSpecs() []definition.EffectSpec {
	return a.effectSpecs
}

func (a Ability) Speed() mtg.Speed {
	return a.speed
}

func (a Ability) Source() query.Object {
	return a.source
}

func (a Ability) Zone() mtg.Zone {
	return a.zone
}

func (a Ability) Name() string {
	//return fmt.Sprintf("%s - %s", a.source.Name(), a.name)
	return a.name
}

func (a Ability) ID() string {
	return a.id
}

// Description returns a string representation of the activated ability.
func (a Ability) Description() string {
	var descriptions []string
	for _, effect := range a.effectSpecs {
		// TODO: Come up with a better way to handle descriptions
		descriptions = append(descriptions, effect.Name)
		//descriptions = append(descriptions, effect.Description())
	}
	// return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(descriptions, ", "))
	return fmt.Sprintf("%s: %s", "<cost>", strings.Join(descriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
// TODO: This needs to happen per effect not per ability.
func (a Ability) IsManaAbility() bool {
	for _, effect := range a.effectSpecs {
		if effect.Name == "AddMana" {
			return true
		}
	}
	return false
}

func (a Ability) Match(predicate query.Predicate) bool {
	return predicate(a)
}
