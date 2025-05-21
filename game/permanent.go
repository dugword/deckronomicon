package game

import (
	"errors"
	"fmt"
)

// Permanent represents a permanent card on the battlefield.
type Permanent struct {
	card               *Card
	tapped             bool
	summoningSickness  bool
	activatedAbilities []*ActivatedAbility
	staticAbilities    []*StaticAbility
	triggeredAbilities []*TriggeredAbility
}

// NewPermanent creates a new Permanent instance from a Card.
func NewPermanent(card *Card) (*Permanent, error) {
	permanent := Permanent{
		card:   card,
		tapped: false,
	}
	if card.HasType(CardTypeCreature) {
		permanent.summoningSickness = true
	}
	var activatedAbilities []*ActivatedAbility
	for _, abilitySpec := range card.object.ActivatedAbilitiesSpec {
		cost, err := NewCost(abilitySpec.CostExpression, &permanent)
		if err != nil {
			return nil, fmt.Errorf("failed to create cost: %w", err)
		}
		var effects []*Effect
		for _, effectSpec := range abilitySpec.EffectSpecs {
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
		activatedAbility := ActivatedAbility{
			Cost:    cost,
			Effects: effects,
		}
		activatedAbilities = append(activatedAbilities, &activatedAbility)
	}
	permanent.activatedAbilities = activatedAbilities
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
	return p.card.object.CardTypes
}

// Colors returns the colors of the permanent.
func (p *Permanent) Colors() *Colors {
	return p.card.object.Colors
}

// Defense returns the defense of the permanent.
func (p *Permanent) Defense() int {
	return p.card.object.Defense
}

// HasType checks if the permanent has the specified card type.
func (p *Permanent) HasType(cardType CardType) bool {
	return p.card.object.HasType(cardType)
}

// ID returns the ID of the permanent.
func (p *Permanent) ID() string {
	return p.card.object.ID
}

// Loyalty returns the loyalty of the permanent.
func (p *Permanent) Loyalty() int {
	return p.card.object.Loyalty
}

// ManCost returns the mana cost of the permanent.
func (p *Permanent) ManaCost() *ManaCost {
	return p.card.object.ManaCost
}

// Name returns the name of the permanent.
func (p *Permanent) Name() string {
	return p.card.object.Name
}

// Power returns the power of the permanent.
func (p *Permanent) Power() int {
	return p.card.object.Power
}

// RulesText returns the rules text of the permanent. The RulesText does not
// impact the game logic.
func (p *Permanent) RulesText() string {
	return p.card.object.RulesText
}

// StaticAbilities returns the static abilities of the permanent.
func (p *Permanent) StaticAbilities() []*StaticAbility {
	return p.staticAbilities
}

// Subtypes returns the subtypes of the permanent.
func (p *Permanent) Subtypes() []Subtype {
	return p.card.object.Subtypes
}

// Supertypes returns the supertypes of the permanent.
func (p *Permanent) Supertypes() []Supertype {
	return p.card.object.Supertypes
}

// Toughness returns the toughness of the permanent.
func (p *Permanent) Toughness() int {
	return p.card.object.Toughness
}

// TriggeredAbilities returns the triggered abilities of the permanent.
func (p *Permanent) TriggeredAbilities() []*TriggeredAbility {
	return p.triggeredAbilities
}
