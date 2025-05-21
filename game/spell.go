package game

import (
	"errors"
	"fmt"
)

// Spell represents a spell object on the stack.
type Spell struct {
	card               *Card
	spellAbility       *SpellAbility
	staticAbilities    []*StaticAbility
	triggeredAbilities []*TriggeredAbility
}

// NewSpell creates a new Spell instance from a Card.
func NewSpell(card *Card) (*Spell, error) {
	spell := Spell{
		card: card,
	}
	// Spell Ability
	if card.SpellAbilitySpec == nil {
		return nil, errors.New("no spell ability")
	}
	// TODO: Additional Costs
	/*
		cost, err := NewCost(card.SpellAbilitySpec.CostExpression, &spell)
		if err != nil {
			return nil, fmt.Errorf("failed to create cost: %w", err)
		}
	*/
	var effects []*Effect
	for _, effectSpec := range card.SpellAbilitySpec.EffectSpecs {
		effectBuilder, ok := EffectMap[effectSpec.ID]
		if !ok {
			return nil, fmt.Errorf("effect %s not found", effectSpec.ID)
		}
		effect, err := effectBuilder(card.Name(), effectSpec.Modifiers)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect %s: %w", effectSpec.ID, err)
		}
		effects = append(effects, effect)
	}
	spellAbility := SpellAbility{
		// Cost:    cost, // TODO: Additional Costs
		Effects: effects,
	}
	// TODO: I think this is wrong, it's overwriting the existing
	// activated abilities of the card. I think we need to shadow this with a
	// top level SpellAbility field.
	spell.spellAbility = &spellAbility
	return &spell, nil
}

// ActivatedAbilities returns the activated abilities of the spell.
// Spells cannot have activated abilities, but this is here for interface
// compatibility.
func (s *Spell) ActivatedAbilities() []ActivatedAbility {
	return []ActivatedAbility{}
}

// Card returns the card associated with the spell.
func (s *Spell) Card() *Card {
	return s.card
}

// CardTypes returns the card types of the spell.
func (s *Spell) CardTypes() []CardType {
	return s.card.object.CardTypes
}

// Colors returns the colors of the spell.
func (s *Spell) Colors() *Colors {
	return s.card.object.Colors
}

// ColorIdicator returns the color indicator of the spell.
func (s *Spell) Defense() int {
	return s.card.object.Defense
}

// HasType checks if the spell has the specified card type.
func (s *Spell) HasType(cardType CardType) bool {
	return s.card.object.HasType(cardType)
}

// ID returns the ID of the spell.
func (s *Spell) ID() string {
	return s.card.object.ID
}

// Loyalty returns the loyalty of the spell.
func (s *Spell) Loyalty() int {
	return s.card.object.Loyalty
}

// ManaCost returns the mana cost of the spell.
func (s *Spell) ManaCost() *ManaCost {
	return s.card.object.ManaCost
}

// Name returns the name of the spell.
func (s *Spell) Name() string {
	return s.card.object.Name
}

// Power returns the power of the spell.
func (s *Spell) Power() int {
	return s.card.object.Power
}

// RulesText returns the rules text of the spell. The RulesText does not
// impact the game logic.
func (s *Spell) RulesText() string {
	return s.card.object.RulesText
}

// SpellAbility returns the spell ability of the spell.
func (s *Spell) SpellAbility() *SpellAbility {
	return s.spellAbility
}

// StaticAbilities returns the static abilities of the spell.
func (s *Spell) StaticAbilities() []*StaticAbility {
	return s.staticAbilities
}

// Subtypes returns the subtypes of the spell.
func (s *Spell) Subtypes() []Subtype {
	return s.card.object.Subtypes
}

// Supertypes returns the supertypes of the spell.
func (s *Spell) Supertypes() []Supertype {
	return s.card.object.Supertypes
}

// Toughness returns the toughness of the spell.
func (s *Spell) Toughness() int {
	return s.card.object.Toughness
}

// TriggeredAbilities returns the triggered abilities of the spell.
func (s *Spell) TriggeredAbilities() []*TriggeredAbility {
	return s.triggeredAbilities
}
