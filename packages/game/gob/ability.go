package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
)

type Ability struct {
	cost    cost.Cost
	effects []effect.Effect
	id      string
	name    string
	source  Object
	speed   mtg.Speed
	zone    mtg.Zone
}

func (a Ability) Controller() string {
	return a.source.Controller()
}

func (a Ability) Description() string {
	return "Put a good description here"
}

func (a Ability) Effects() []effect.Effect {
	return a.effects
}

func (a Ability) Cost() cost.Cost {
	return a.cost
}

func (a Ability) ID() string {
	return a.id
}

func (a Ability) Match(predicate query.Predicate) bool {
	return predicate(a)
}

func (a Ability) Name() string {
	return a.name
}

func (a Ability) Owner() string {
	return a.source.Owner()
}

func (a Ability) Speed() mtg.Speed {
	return a.speed
}

func (a Ability) Source() Object {
	return a.source
}

func (a Ability) Zone() mtg.Zone {
	return a.zone
}

func buildAbility(source Object, idx int, abilityDefinition definition.Ability) (Ability, error) {
	speed := mtg.SpeedInstant
	specSpeed, ok := mtg.StringToSpeed(abilityDefinition.Speed)
	if ok {
		speed = specSpeed
	}
	zone := mtg.ZoneBattlefield
	specZone, ok := mtg.StringToZone(abilityDefinition.Zone)
	if ok {
		zone = specZone
	}
	abilityCost, err := cost.Parse(abilityDefinition.Cost)
	if err != nil {
		return Ability{}, fmt.Errorf("failed to parse cost %q: %w", abilityDefinition.Cost, err)
	}
	var effects []effect.Effect
	for _, effectDefinition := range abilityDefinition.Effects {
		effect, err := effect.New(effectDefinition)
		if err != nil {
			return Ability{}, fmt.Errorf("failed to build effect %s: %w", effectDefinition.Name, err)
		}
		effects = append(effects, effect)
	}
	ability := Ability{
		cost:    abilityCost,
		name:    abilityDefinition.Name,
		effects: effects,
		id:      fmt.Sprintf("%s-%d", source.ID(), idx+1),
		zone:    zone,
		speed:   speed,
		source:  source,
	}
	return ability, nil
}
