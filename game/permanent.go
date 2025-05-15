package game

import "errors"

type Permanent struct {
	card              *Card
	tapped            bool
	summoningSickness bool
}

func NewPermanent(card *Card) *Permanent {
	permanent := Permanent{
		card:   card,
		tapped: false,
	}
	if card.HasType(CardTypeCreature) {
		permanent.summoningSickness = true
	}

	return &permanent
}

func (p *Permanent) IsTapped() bool {
	return p.tapped
}

func (p *Permanent) Tap() error {
	if p.tapped {
		return errors.New("already tapped")
	}
	p.tapped = true
	return nil
}

func (p *Permanent) Untap() {
	p.tapped = false
}

func (p *Permanent) HasSummoningSickness() bool {
	return p.summoningSickness
}

func (p *Permanent) RemoveSummoningSickness() {
	p.summoningSickness = false
}

func (p *Permanent) ActivatedAbilities() []ActivatedAbility {
	return p.card.object.ActivatedAbilities
}
func (p *Permanent) Card() *Card {
	return p.card
}
func (p *Permanent) CardTypes() []CardType {
	return p.card.object.CardTypes
}
func (p *Permanent) Color() string {
	return p.card.object.Color
}
func (p *Permanent) ColorIdicator() string {
	return p.card.object.ColorIdicator
}
func (p *Permanent) Defense() int {
	return p.card.object.Defense
}
func (p *Permanent) HasType(cardType CardType) bool {
	return p.card.object.HasType(cardType)
}
func (p *Permanent) ID() string {
	return p.card.object.ID
}
func (p *Permanent) Loyalty() int {
	return p.card.object.Loyalty
}
func (p *Permanent) ManaCost() ManaCost {
	return p.card.object.ManaCost
}
func (p *Permanent) Name() string {
	return p.card.object.Name
}
func (p *Permanent) Power() int {
	return p.card.object.Power
}
func (p *Permanent) RulesText() string {
	// TODO: Need to build this from abilities....
	return p.card.object.RulesText
}
func (p *Permanent) SpellAbility() *SpellAbility {
	return p.card.object.SpellAbility
}
func (p *Permanent) StaticAbilities() []StaticAbility {
	return p.card.object.StaticAbilities
}
func (p *Permanent) Subtypes() []Subtype {
	return p.card.object.Subtypes
}
func (p *Permanent) Supertypes() []Supertype {
	return p.card.object.Supertypes
}
func (p *Permanent) Toughness() int {
	return p.card.object.Toughness
}
func (p *Permanent) TriggeredAbilities() []TriggeredAbility {
	return p.card.object.TriggeredAbilities
}
