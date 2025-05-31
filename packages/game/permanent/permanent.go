package permanent

import (
	"deckronomicon/packages/game/ability/activated"
	"deckronomicon/packages/game/ability/static"
	"deckronomicon/packages/game/ability/triggered"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
)

type State interface {
	GetNextID() string
}

// Permanent represents a permanent card on the battlefield.
type Permanent struct {
	activatedAbilities []*activated.Ability
	card               *card.Card
	cardTypes          []mtg.CardType
	colors             mtg.Colors
	id                 string
	loyalty            int
	manaCost           *cost.ManaCost
	name               string
	power              int
	rulesText          string
	staticAbilities    []*static.Ability
	subtypes           []mtg.Subtype
	summoningSickness  bool
	supertypes         []mtg.Supertype
	tapped             bool
	toughness          int
	triggeredAbilities []*triggered.Ability
}

// NewPermanent creates a new Permanent instance from a Card.
func NewPermanent(card *card.Card, state State) (*Permanent, error) {
	permanent := Permanent{
		activatedAbilities: []*activated.Ability{},
		card:               card,
		cardTypes:          card.CardTypes(),
		colors:             card.Colors(),
		id:                 state.GetNextID(),
		loyalty:            card.Loyalty(),
		manaCost:           card.ManaCost(),
		name:               card.Name(),
		power:              card.Power(),
		rulesText:          card.RulesText(),
		staticAbilities:    card.StaticAbilities(),
		subtypes:           card.Subtypes(),
		supertypes:         card.Supertypes(),
		toughness:          card.Toughness(),
		triggeredAbilities: []*triggered.Ability{},
	}
	if card.Match(has.CardType(mtg.CardTypeCreature)) {
		permanent.summoningSickness = true
	}
	return &permanent, nil
}

// ActivatedAbilities returns the activated abilities of the permanent.
func (p *Permanent) ActivatedAbilities() []*activated.Ability {
	return p.activatedAbilities
}

// Card returns the card associated with the permanent.
func (p *Permanent) Card() *card.Card {
	return p.card
}

// CardTypes returns the card types of the permanent.
func (p *Permanent) CardTypes() []mtg.CardType {
	return p.cardTypes
}

// Colors returns the colors of the permanent.
func (p *Permanent) Colors() mtg.Colors {
	return p.colors
}

func (p *Permanent) ID() string {
	return p.id
}

// IsTapped checks if the permanent is tapped.
func (p *Permanent) IsTapped() bool {
	return p.tapped
}

// HasSummoningSickness checks if the permanent has summoning sickness.
func (p *Permanent) HasSummoningSickness() bool {
	return p.summoningSickness
}

// Loyalty returns the loyalty of the permanent.
func (p *Permanent) Loyalty() int {
	return p.loyalty
}

// ManCost returns the mana cost of the permanent.
func (p *Permanent) ManaCost() *cost.ManaCost {
	return p.manaCost
}

func (p *Permanent) ManaValue() int {
	if p.manaCost == nil {
		return 0
	}
	return p.manaCost.ManaValue()
}

func (per *Permanent) Match(p query.Predicate) bool {
	return p(per)
}

// Name returns the name of the permanent.
func (p *Permanent) Name() string {
	return p.name
}

// Power returns the power of the permanent.
func (p *Permanent) Power() int {
	return p.power
}

// RemoveSummoningSickness removes summoning sickness from the permanent. It
// is a valid operation even if the permanent does not have summoning
// sickness.
func (p *Permanent) RemoveSummoningSickness() {
	p.summoningSickness = false
}

// RulesText returns the rules text of the permanent. The RulesText does not
// impact the game logic.
func (p *Permanent) RulesText() string {
	return p.rulesText
}

// StaticAbilities returns the static abilities of the permanent.
func (p *Permanent) StaticAbilities() []*static.Ability {
	return p.staticAbilities
}

// Subtypes returns the subtypes of the permanent.
func (p *Permanent) Subtypes() []mtg.Subtype {
	return p.subtypes
}

// Supertypes returns the supertypes of the permanent.
func (p *Permanent) Supertypes() []mtg.Supertype {
	return p.supertypes
}

// Tap taps the permanent. Returns an error if the permanent is already
// tapped.
// TODO: There are some specifics that matter for tapping tapped creatures,
// but for now we will just disallow it.
func (p *Permanent) Tap() error {
	if p.tapped {
		return mtg.ErrAlreadyTapped
	}
	p.tapped = true
	return nil
}

// Toughness returns the toughness of the permanent.
func (p *Permanent) Toughness() int {
	return p.toughness
}

// TriggeredAbilities returns the triggered abilities of the permanent.
func (p *Permanent) TriggeredAbilities() []*triggered.Ability {
	return p.triggeredAbilities
}

// Untap untaps the permanent. It is a valid operation even if the permanent
// is already untapped.
// TODO: There are some specifics that matter for untapping untapped
// creatures, but for now we will just allow it.
func (p *Permanent) Untap() {
	p.tapped = false
}
