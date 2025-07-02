package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/query"
	"fmt"
)

type Spell struct {
	card              *Card
	cardTypes         []mtg.CardType
	colors            mtg.Colors
	controller        string
	effectWithTargets []*effect.EffectWithTarget
	flashback         bool
	id                string
	isCopy            bool
	loyalty           int
	manaCost          cost.Mana
	name              string
	owner             string
	power             int
	rulesText         string
	staticAbilities   []staticability.StaticAbility
	subtypes          []mtg.Subtype
	supertypes        []mtg.Supertype
	toughness         int
}

func CopySpell(
	id string,
	spell *Spell,
	playerID string,
	effectWithTargets []*effect.EffectWithTarget,
) (*Spell, error) {
	copiedSpell, err := NewSpell(
		id,
		spell.card,
		playerID,
		effectWithTargets,
		spell.flashback,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create copied spell: %w", err)
	}
	copiedSpell.isCopy = true
	return copiedSpell, nil
}

func NewSpell(
	id string,
	card *Card,
	playerID string,
	effectWithTargets []*effect.EffectWithTarget,
	flashback bool,
) (*Spell, error) {
	spell := Spell{
		card:              card,
		cardTypes:         card.CardTypes(),
		colors:            card.Colors(),
		controller:        playerID,
		effectWithTargets: effectWithTargets,
		flashback:         flashback,
		id:                id,
		loyalty:           card.Loyalty(),
		manaCost:          card.ManaCost(),
		name:              card.Name(),
		owner:             card.Owner(),
		power:             card.Power(),
		rulesText:         card.RulesText(),
		staticAbilities:   card.StaticAbilities(),
		subtypes:          card.Subtypes(),
		supertypes:        card.Supertypes(),
		toughness:         card.Toughness(),
	}
	return &spell, nil
}

func (s *Spell) Card() *Card {
	return s.card
}

func (s *Spell) CardTypes() []mtg.CardType {
	return s.cardTypes
}

func (s *Spell) Colors() mtg.Colors {
	return s.colors
}

func (s *Spell) Controller() string {
	return s.controller
}

func (s *Spell) EffectWithTargets() []*effect.EffectWithTarget {
	return s.effectWithTargets
}

func (s *Spell) Flashback() bool {
	return s.flashback
}

func (s *Spell) Description() string {
	return "Put something good here"
}

func (s *Spell) ID() string {
	return s.id
}

func (s *Spell) IsCopy() bool {
	return s.isCopy
}

func (s *Spell) Loyalty() int {
	return s.loyalty
}

func (s *Spell) ManaCost() cost.Mana {
	return s.manaCost
}

func (s *Spell) ManaValue() int {
	return s.manaCost.Amount().Total()
}

func (s *Spell) Match(predicate query.Predicate) bool {
	return predicate(s)
}

func (s *Spell) Name() string {
	return s.name
}

func (s *Spell) Owner() string {
	return s.owner
}

func (s *Spell) Power() int {
	return s.power
}

func (s *Spell) RulesText() string {
	return s.rulesText
}

func (s *Spell) SourceID() string {
	return s.card.id
}

func (s *Spell) StaticAbilities() []staticability.StaticAbility {
	return s.staticAbilities
}

func (s *Spell) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range s.staticAbilities {
		if ability.StaticKeyword() != "" {
			keywords = append(keywords, ability.StaticKeyword())
		}
	}
	return keywords
}

func (s *Spell) Subtypes() []mtg.Subtype {
	return s.subtypes
}

func (s *Spell) Supertypes() []mtg.Supertype {
	return s.supertypes
}

func (s *Spell) Toughness() int {
	return s.toughness
}
