package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
)

// Spell represents a spell object on the stack.
type Spell struct {
	card              Card
	cardTypes         []mtg.CardType
	colors            mtg.Colors
	controller        string
	owner             string
	id                string
	loyalty           int
	manaCost          cost.ManaCost
	name              string
	power             int
	rulesText         string
	effectWithTargets []EffectWithTarget
	staticAbilities   []StaticAbility
	subtypes          []mtg.Subtype
	supertypes        []mtg.Supertype
	toughness         int
	flashback         bool
	isCopy            bool
}

func CopySpell(
	id string,
	spell Spell,
	playerID string,
	effectWithTargets []EffectWithTarget,
) (Spell, error) {
	copiedSpell, err := NewSpell(
		id,
		spell.card,
		playerID,
		effectWithTargets,
		spell.flashback,
	)
	if err != nil {
		return Spell{}, fmt.Errorf("failed to create copied spell: %w", err)
	}
	copiedSpell.isCopy = true
	return copiedSpell, nil
}

// NewSpell creates a new Spell instance from a Card.
func NewSpell(
	id string,
	card Card,
	playerID string,
	effectWithTargets []EffectWithTarget,
	flashback bool,
) (Spell, error) {
	spell := Spell{
		card:              card,
		cardTypes:         card.CardTypes(),
		colors:            card.Colors(),
		controller:        playerID,
		effectWithTargets: effectWithTargets,
		flashback:         flashback,
		owner:             card.Owner(),
		id:                id,
		loyalty:           card.Loyalty(),
		manaCost:          card.ManaCost(),
		name:              card.Name(),
		power:             card.Power(),
		rulesText:         card.RulesText(),
		staticAbilities:   card.StaticAbilities(),
		subtypes:          card.Subtypes(),
		supertypes:        card.Supertypes(),
		toughness:         card.Toughness(),
	}
	return spell, nil
}

func (s Spell) Card() Card {
	return s.card
}

// CardTypes returns the card types of the spell.
func (s Spell) CardTypes() []mtg.CardType {
	return s.cardTypes
}

// Colors returns the colors of the spell.
func (s Spell) Colors() mtg.Colors {
	return s.colors
}

func (s Spell) Controller() string {
	return s.controller
}

// Effects returns the effects of the spell.
func (s Spell) EffectWithTargets() []EffectWithTarget {
	return s.effectWithTargets
}

func (s Spell) Flashback() bool {
	return s.flashback
}

// Description returns a string representation of the activated ability.
func (s Spell) Description() string {
	return "Put something good here"
}

// ID returns the ID of the spell.
func (s Spell) ID() string {
	return s.id
}

// Loyalty returns the loyalty of the spell.
func (s Spell) Loyalty() int {
	return s.loyalty
}

// ManaCost returns the mana cost of the spell.
func (s Spell) ManaCost() cost.ManaCost {
	return s.manaCost
}

func (s Spell) ManaValue() int {
	return s.manaCost.Amount().Total()
}

func (s Spell) Match(predicate query.Predicate) bool {
	return predicate(s)
}

// Name returns the name of the spell.
func (s Spell) Name() string {
	return s.name
}

func (s Spell) Owner() string {
	return s.owner
}

// Power returns the power of the spell.
func (s Spell) Power() int {
	return s.power
}

// RulesText returns the rules text of the spell. The RulesText does not
// impact the game logic.
func (s Spell) RulesText() string {
	return s.rulesText
}

func (s Spell) SourceID() string {
	return s.card.id
}

// StaticAbilities returns the static abilities of the spell
func (s Spell) StaticAbilities() []StaticAbility {
	return s.staticAbilities
}

func (s Spell) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range s.staticAbilities {
		keyword, ok := mtg.StringToStaticKeyword(ability.Name())
		if !ok {
			continue
		}
		keywords = append(keywords, keyword)
	}
	return keywords
}

// Subtypes returns the subtypes of the spell.
func (s Spell) Subtypes() []mtg.Subtype {
	return s.subtypes
}

// Supertypes returns the supertypes of the spell.
func (s Spell) Supertypes() []mtg.Supertype {
	return s.supertypes
}

// Toughness returns the toughness of the spell.
func (s Spell) Toughness() int {
	return s.toughness
}

func (s Spell) IsCopy() bool {
	return s.isCopy
}
