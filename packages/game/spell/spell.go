package spell

import (
	"deckronomicon/packages/game/ability/spell"
	"deckronomicon/packages/game/ability/static"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"deckronomicon/packages/query/predicate"
	"fmt"
	"strings"
)

type Battlefield interface {
	Add(permanent *permanent.Permanent) error
}

// This sucks that this is here, it's only to satisfy the interface. Maybe
// call report state in the caller instead.
type Agent interface {
	ReportState(State)
}

// Spell represents a spell object on the stack.
type Spell struct {
	card            *card.Card
	cardTypes       []mtg.CardType
	colors          mtg.Colors
	flashback       bool
	id              string
	loyalty         int
	manaCost        *cost.ManaCost
	name            string
	power           int
	rulesText       string
	spellAbility    *spell.Ability
	staticAbilities []*static.Ability
	subtypes        []mtg.Subtype
	supertypes      []mtg.Supertype
	toughness       int
}

// NewSpell creates a new Spell instance from a Card.
func NewSpell(state state, card *card.Card) (*Spell, error) {
	spell := Spell{
		card:            card,
		cardTypes:       card.CardTypes(),
		colors:          card.Colors(),
		id:              state.GetNextID(),
		loyalty:         card.Loyalty(),
		manaCost:        card.ManaCost(),
		name:            card.Name(),
		power:           card.Power(),
		rulesText:       card.RulesText(),
		staticAbilities: card.StaticAbilities(),
		subtypes:        card.Subtypes(),
		supertypes:      card.Supertypes(),
		toughness:       card.Toughness(),
	}
	return &spell, nil
}

func (s *Spell) Card() *card.Card {
	return s.card
}

// CardTypes returns the card types of the spell.
func (s *Spell) CardTypes() []mtg.CardType {
	return s.cardTypes
}

// Colors returns the colors of the spell.
func (s *Spell) Colors() mtg.Colors {
	return s.colors
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
func (s *Spell) ManaCost() *cost.ManaCost {
	return s.manaCost
}

func (s *Spell) ManaValue() int {
	if s.manaCost == nil {
		return 0
	}
	return s.manaCost.ManaValue()
}

func (s *Spell) Match(p predicate.Predicate) bool {
	return p(s)
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
		descriptions = append(descriptions, effect.Description())
	}
	// TODO: Support additional costs
	return fmt.Sprintf("%s: %s", s.ManaCost().Description(), strings.Join(descriptions, ", "))
}

/*
func (s *Spell) Resolve(state state, player player) error {
	if s.spellAbility == nil && s.card.Match(is.Permanent()) {
		permanent, err := permanent.NewPermanent(s.card, state)
		if err != nil {
			return fmt.Errorf("failed to create permanent: %w", err)
		}
		player.Battlefield().Add(permanent)
	}
		if err := s.spellAbility.Resolve(state, player); err != nil {
			return fmt.Errorf("cannot resolve spell ability: %w", err)
		}
	if s.flashback {
		if err := player.Exile.Add(s.card); err != nil {
			return fmt.Errorf("cannot move spell to exile: %w", err)
		}
		return nil
	}
	if err := player.Graveyard.Add(s.card); err != nil {
		return fmt.Errorf("cannot move spell to graveyard: %w", err)
	}
	return nil
}
*/

// RulesText returns the rules text of the spell. The RulesText does not
// impact the game logic.
func (s *Spell) RulesText() string {
	return s.rulesText
}

// SpellAbility returns the spell ability of the spell.
func (s *Spell) SpellAbility() *spell.Ability {
	return s.spellAbility
}

// StaticAbilities returns the static abilities of the spell
func (s *Spell) StaticAbilities() []*static.Ability {
	return s.staticAbilities
}

// Subtypes returns the subtypes of the spell.
func (s *Spell) Subtypes() []mtg.Subtype {
	return s.subtypes
}

// Supertypes returns the supertypes of the spell.
func (s *Spell) Supertypes() []mtg.Supertype {
	return s.supertypes
}

// Toughness returns the toughness of the spell.
func (s *Spell) Toughness() int {
	return s.toughness
}

/*
func (s *Spell) Flashback() {
	s.flashback = true
}

func (s *Spell) Splice(state state, card *card.Card) error {
	// TODO: This is what was missing
	spell, err := NewSpell(state, card)
	if err != nil {
		return fmt.Errorf("failed to create spell for splice: %w", err)
	}
	s.spellAbility.Splice(spell.spellAbility)
	return nil
}
*/
