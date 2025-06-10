package gob

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

type CardInZone struct {
	card Card
	zone mtg.Zone
}

type AbilityInZone struct {
	zone    mtg.Zone
	ability Ability
}

func NewAbilityInZone(ability Ability, object query.Object, zone mtg.Zone) AbilityInZone {
	return AbilityInZone{
		zone:    zone,
		ability: ability,
	}
}

func (a AbilityInZone) Match(predicate query.Predicate) bool {
	return predicate(a.ability)
}

func (a AbilityInZone) Description() string {
	return a.ability.Description()
}

func (a AbilityInZone) Name() string {
	return a.ability.Name()
}

func (a AbilityInZone) ID() string {
	return a.ability.ID()
}

func (a AbilityInZone) Zone() mtg.Zone {
	return a.zone
}

func (a AbilityInZone) Ability() Ability {
	return a.ability
}

func (a AbilityInZone) Source() query.Object {
	return a.ability.Source()
}

func NewCardInZone(card Card, zone mtg.Zone) CardInZone {
	return CardInZone{
		card: card,
		zone: zone,
	}
}

func (c CardInZone) Description() string {
	return c.card.Description()
}

func (c CardInZone) ID() string {
	return c.card.ID()
}

func (c CardInZone) Match(predicate query.Predicate) bool {
	return predicate(c.card)
}

func (c CardInZone) Name() string {
	return c.card.Name()
}

func (c CardInZone) Zone() mtg.Zone {
	return c.zone
}

func (c CardInZone) Card() Card {
	return c.card
}

type ObjectInZone struct {
	object query.Object
	zone   mtg.Zone
}

func NewObjectInZone(object query.Object, zone mtg.Zone) ObjectInZone {
	return ObjectInZone{
		object: object,
		zone:   zone,
	}
}

func (o ObjectInZone) Name() string {
	return o.object.Name()
}

func (o ObjectInZone) ID() string {
	return o.object.ID()
}

func (o ObjectInZone) Zone() mtg.Zone {
	return o.zone
}

func (o ObjectInZone) Card() Card {
	if card, ok := o.object.(Card); ok {
		return card
	}
	panic("ObjectInZone does not contain a Card")
}
