package game

import (
	"errors"
	"fmt"
)

// Permanent represents a permanent card on the battlefield.
type Permanent struct {
	activatedAbilities []*ActivatedAbility
	card               *Card
	cardTypes          []CardType
	colors             *Colors
	id                 string
	loyalty            int
	manaCost           *ManaCost
	name               string
	power              int
	rulesText          string
	staticAbilities    []*StaticAbility
	subtypes           []Subtype
	summoningSickness  bool
	supertypes         []Supertype
	tapped             bool
	toughness          int
	triggeredAbilities []*TriggeredAbility
}

// NewPermanent creates a new Permanent instance from a Card.
func NewPermanent(card *Card) (*Permanent, error) {
	permanent := Permanent{
		activatedAbilities: []*ActivatedAbility{},
		card:               card,
		cardTypes:          card.cardTypes,
		colors:             card.colors,
		id:                 GetNextID(),
		loyalty:            card.loyalty,
		manaCost:           card.manaCost,
		name:               card.name,
		power:              card.power,
		rulesText:          card.rulesText,
		staticAbilities:    []*StaticAbility{},
		subtypes:           card.subtypes,
		supertypes:         card.supertypes,
		toughness:          card.toughness,
		triggeredAbilities: []*TriggeredAbility{},
	}
	if card.HasCardType(CardTypeCreature) {
		permanent.summoningSickness = true
	}
	for _, spec := range card.activatedAbilitySpecs {
		if spec.Zone == ZoneHand || spec.Zone == ZoneGraveyard {
			continue
		}
		ability, err := BuildActivatedAbility(*spec, &permanent)
		if err != nil {
			return nil, fmt.Errorf("failed to build activated ability: %w", err)
		}
		permanent.activatedAbilities = append(permanent.activatedAbilities, ability)
	}
	for _, spec := range card.staticAbilitySpecs {
		if spec.Zone == ZoneHand || spec.Zone == ZoneGraveyard {
			continue
		}
		ability, err := BuildStaticAbility(*spec, &permanent)
		if err != nil {
			return nil, fmt.Errorf("failed to build activated ability: %w", err)
		}
		permanent.staticAbilities = append(permanent.staticAbilities, ability)
	}
	return &permanent, nil
}

// IsTapped checks if the permanent is tapped.
func (p *Permanent) IsTapped() bool {
	return p.tapped
}

// Tap taps the permanent. Returns an error if the permanent is already
// tapped.
func (p *Permanent) Tap() error {
	if p.tapped {
		return errors.New("already tapped")
	}
	p.tapped = true
	return nil
}

// Untap untaps the permanent. It is a valid operation even if the permanent
// is already untapped.
func (p *Permanent) Untap() {
	p.tapped = false
}

// HasSummoningSickness checks if the permanent has summoning sickness.
func (p *Permanent) HasSummoningSickness() bool {
	return p.summoningSickness
}

// HasSubtype checks if the card has a specific type. It returns true if the
// card has the specified type, and false otherwise.
func (p *Permanent) HasSubtype(subtype Subtype) bool {
	for _, t := range p.subtypes {
		if t == subtype {
			return true
		}
	}
	return false
}

// RemoveSummoningSickness removes summoning sickness from the permanent. It
// is a valid operation even if the permanent does not have summoning
// sickness.
func (p *Permanent) RemoveSummoningSickness() {
	p.summoningSickness = false
}

// ActivatedAbilities returns the activated abilities of the permanent.
func (p *Permanent) ActivatedAbilities() []*ActivatedAbility {
	return p.activatedAbilities
}

// Card returns the card associated with the permanent.
func (p *Permanent) Card() *Card {
	return p.card
}

// CardTypes returns the card types of the permanent.
func (p *Permanent) CardTypes() []CardType {
	return p.cardTypes
}

// Colors returns the colors of the permanent.
func (p *Permanent) Colors() *Colors {
	return p.colors
}

// HasType checks if the permanent has the specified card type.
func (p *Permanent) HasCardType(cardType CardType) bool {
	for _, t := range p.cardTypes {
		if t == cardType {
			return true
		}
	}
	return false
}

// ID returns the ID of the permanent.
func (p *Permanent) ID() string {
	return p.id
}

// Loyalty returns the loyalty of the permanent.
func (p *Permanent) Loyalty() int {
	return p.loyalty
}

// ManCost returns the mana cost of the permanent.
func (p *Permanent) ManaCost() *ManaCost {
	return p.manaCost
}

// Name returns the name of the permanent.
func (p *Permanent) Name() string {
	return p.name
}

// Power returns the power of the permanent.
func (p *Permanent) Power() int {
	return p.power
}

// RulesText returns the rules text of the permanent. The RulesText does not
// impact the game logic.
func (p *Permanent) RulesText() string {
	return p.rulesText
}

// StaticAbilities returns the static abilities of the permanent.
func (p *Permanent) StaticAbilities() []*StaticAbility {
	return p.staticAbilities
}

// Subtypes returns the subtypes of the permanent.
func (p *Permanent) Subtypes() []Subtype {
	return p.subtypes
}

// Supertypes returns the supertypes of the permanent.
func (p *Permanent) Supertypes() []Supertype {
	return p.supertypes
}

// Toughness returns the toughness of the permanent.
func (p *Permanent) Toughness() int {
	return p.toughness
}

// TriggeredAbilities returns the triggered abilities of the permanent.
func (p *Permanent) TriggeredAbilities() []*TriggeredAbility {
	return p.triggeredAbilities
}
