package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/query"
)

type Card struct {
	activatedAbilities []Ability
	cardTypes          []mtg.CardType
	controller         string
	owner              string
	colors             mtg.Colors
	id                 string
	loyalty            int
	manaCost           cost.Mana
	additionalCost     cost.Cost
	name               string
	power              int
	rulesText          string
	spellAbility       []effect.Effect
	staticAbilities    []staticability.StaticAbility
	subtypes           []mtg.Subtype
	supertypes         []mtg.Supertype
	toughness          int
}

func (c Card) ActivatedAbilities() []Ability {
	return c.activatedAbilities[:]
}

func (c Card) AdditionalCost() cost.Cost {
	return c.additionalCost
}

func (c Card) CardTypes() []mtg.CardType {
	return c.cardTypes
}

func (c Card) Colors() mtg.Colors {
	return c.colors
}

func (c Card) Controller() string {
	return c.controller
}

func (c Card) Description() string {
	return c.rulesText
}

func (c Card) ID() string {
	return c.id
}

func (c Card) Match(predicate query.Predicate) bool {
	return predicate(c)
}

func (c Card) Loyalty() int {
	return c.loyalty
}

func (c Card) ManaCost() cost.Mana {
	return c.manaCost
}

func (c Card) ManaValue() int {
	return c.manaCost.Amount().Total()
}

func (c Card) Name() string {
	return c.name
}

func (c Card) Owner() string {
	return c.owner
}

func (c Card) Power() int {
	return c.power
}
func (c Card) RulesText() string {
	return c.rulesText
}

func (c Card) SpellAbility() []effect.Effect {
	return c.spellAbility[:]
}

func (c Card) StaticAbilities() []staticability.StaticAbility {
	return c.staticAbilities[:]
}

func (c Card) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range c.staticAbilities {
		if ability.StaticKeyword() != "" {
			keywords = append(keywords, ability.StaticKeyword())
		}
	}
	return keywords
}

func (c Card) StaticAbility(keyword mtg.StaticKeyword) (staticability.StaticAbility, bool) {
	for _, ability := range c.staticAbilities {
		if ability.StaticKeyword() == keyword {
			return ability, true
		}
	}
	return nil, false
}

func (c Card) Subtypes() []mtg.Subtype {
	return c.subtypes
}

func (c Card) Supertypes() []mtg.Supertype {
	return c.supertypes
}

func (c Card) Toughness() int {
	return c.toughness
}
