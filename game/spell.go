package game

import (
	"fmt"
	"strings"
)

// Spell represents a spell object on the stack.
type Spell struct {
	card            *Card
	cardTypes       []CardType
	colors          *Colors
	id              string
	loyalty         int
	manaCost        *ManaCost
	name            string
	power           int
	rulesText       string
	spellAbility    *SpellAbility
	staticAbilities []*StaticAbility
	subtypes        []Subtype
	supertypes      []Supertype
	toughness       int
}

// NewSpell creates a new Spell instance from a Card.
func NewSpell(card *Card) (*Spell, error) {
	spell := Spell{
		card:            card,
		cardTypes:       card.cardTypes,
		colors:          card.colors,
		id:              GetNextID(),
		loyalty:         card.loyalty,
		manaCost:        card.manaCost,
		name:            card.name,
		power:           card.power,
		rulesText:       card.rulesText,
		staticAbilities: card.staticAbilities,
		subtypes:        card.subtypes,
		supertypes:      card.supertypes,
		toughness:       card.toughness,
	}
	// TODO: Additional Costs
	/*
		cost, err := NewCost(card.SpellAbilitySpec.CostExpression, &spell)
		if err != nil {
			return nil, fmt.Errorf("failed to create cost: %w", err)
		}
	*/
	spellAbility, err := BuildSpellAbility(card.spellAbilitySpec, &spell)
	if err != nil {
		return nil, fmt.Errorf("failed to build spell ability: %w", err)
	}
	spell.spellAbility = spellAbility
	return &spell, nil
}

// ActivatedAbilities don't exist for spell objects. This method is here to
// satisfy the GameObject interface.
func (s *Spell) ActivatedAbilities() []*ActivatedAbility {
	return nil
}

func (s *Spell) Card() *Card {
	return s.card
}

// CardTypes returns the card types of the spell.
func (s *Spell) CardTypes() []CardType {
	return s.cardTypes
}

// Colors returns the colors of the spell.
func (s *Spell) Colors() *Colors {
	return s.colors
}

// HasCardType checks if the spell has the specified card type.
func (s *Spell) HasCardType(cardType CardType) bool {
	for _, t := range s.cardTypes {
		if t == cardType {
			return true
		}
	}
	return false
}

// HasStaticAbility checks if the spell has a static ability with the given
// ID.
func (s *Spell) HasStaticAbility(id string) bool {
	for _, ability := range s.staticAbilities {
		if ability.ID == id {
			return true
		}
	}
	return false
}

// HasSubtype checks if the card has a specific type. It returns true if the
// card has the specified type, and false otherwise.
func (s *Spell) HasSubtype(subtype Subtype) bool {
	for _, t := range s.subtypes {
		if t == subtype {
			return true
		}
	}
	return false
}

// HasType checks if the spell has the specified card type.
func (s *Spell) HasType(cardType CardType) bool {
	return s.HasType(cardType)
}

// ID returns the ID of the spell.
func (s *Spell) ID() string {
	return s.id
}

// Loyalty returns the loyalty of the spell.
func (s *Spell) Loyalty() int {
	return s.loyalty
}

// ManaCost returns the mana cost of the spell.
func (s *Spell) ManaCost() *ManaCost {
	return s.manaCost
}

func (s *Spell) ManaValue() int {
	if s.manaCost == nil {
		return 0
	}
	return s.manaCost.ManaValue()
}

// Name returns the name of the spell.
func (s *Spell) Name() string {
	return s.name
}

// Power returns the power of the spell.
func (s *Spell) Power() int {
	return s.power
}

// Description returns a string representation of the activated ability.
func (s *Spell) Description() string {
	var descriptions []string
	for _, effect := range s.spellAbility.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	// TODO: Support additional costs
	return fmt.Sprintf("%s: %s", s.ManaCost().Description(), strings.Join(descriptions, ", "))
}

func (s *Spell) Resolve(state *GameState, player *Player) error {
	if s.spellAbility == nil && s.card.IsPermanent() {
		permanent, err := NewPermanent(s.card)
		if err != nil {
			return fmt.Errorf("failed to create permanent: %w", err)
		}
		player.Battlefield.Add(permanent)
	}
	if err := s.spellAbility.Resolve(state, player); err != nil {
		return fmt.Errorf("cannot resolve spell ability: %w", err)
	}
	if err := player.Graveyard.Add(s.card); err != nil {
		return fmt.Errorf("cannot move spell to graveyard: %w", err)
	}
	return nil
}

// RulesText returns the rules text of the spell. The RulesText does not
// impact the game logic.
func (s *Spell) RulesText() string {
	return s.rulesText
}

// SpellAbility returns the spell ability of the spell.
func (s *Spell) SpellAbility() *SpellAbility {
	return s.spellAbility
}

// StaticAbilities returns the static abilities of the spell
func (s *Spell) StaticAbilities() []*StaticAbility {
	return s.staticAbilities
}

// Subtypes returns the subtypes of the spell.
func (s *Spell) Subtypes() []Subtype {
	return s.subtypes
}

// Supertypes returns the supertypes of the spell.
func (s *Spell) Supertypes() []Supertype {
	return s.supertypes
}

// Toughness returns the toughness of the spell.
func (s *Spell) Toughness() int {
	return s.toughness
}
