package gob

import (

	// "deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/query"
	"fmt"
	"slices"
)

type Permanent struct {
	activatedAbilities []Ability
	card               Card
	cardTypes          []mtg.CardType
	colors             mtg.Colors
	controller         string
	id                 string
	loyalty            int
	manaCost           cost.ManaCost
	name               string
	owner              string
	power              int
	rulesText          string
	staticAbilities    []staticability.StaticAbility
	subtypes           []mtg.Subtype
	summoningSickness  bool
	supertypes         []mtg.Supertype
	tapped             bool
	toughness          int
	triggeredAbilities []Ability
}

func NewPermanent(id string, card Card, playerID string) (Permanent, error) {
	permanent := Permanent{
		card:               card,
		cardTypes:          card.CardTypes(),
		colors:             card.Colors(),
		controller:         playerID,
		id:                 id,
		loyalty:            card.Loyalty(),
		manaCost:           card.ManaCost(),
		name:               card.Name(),
		owner:              playerID,
		power:              card.Power(),
		rulesText:          card.RulesText(),
		staticAbilities:    card.StaticAbilities(),
		subtypes:           card.Subtypes(),
		supertypes:         card.Supertypes(),
		toughness:          card.Toughness(),
		triggeredAbilities: []Ability{},
	}
	if slices.Contains(permanent.cardTypes, mtg.CardTypeCreature) {
		permanent.summoningSickness = true
	}
	for i, a := range card.ActivatedAbilities() {
		ability := Ability{
			effects: a.Effects(),
			name:    a.Name(),
			cost:    a.Cost(),
			id:      fmt.Sprintf("%s-%d", id, i+1),
			zone:    a.zone,
			speed:   a.speed,
			source:  permanent,
		}
		permanent.activatedAbilities = append(permanent.activatedAbilities, ability)
	}
	return permanent, nil
}

func (p Permanent) ActivatedAbilities() []Ability {
	return p.activatedAbilities
}

func (p Permanent) Card() Card {
	return p.card
}
func (p Permanent) CardTypes() []mtg.CardType {
	return p.cardTypes
}

func (p Permanent) Colors() mtg.Colors {
	return p.colors
}

func (p Permanent) Controller() string {
	return p.controller
}

func (p Permanent) Description() string {
	return p.rulesText
}

func (p Permanent) ID() string {
	return p.id
}

func (p Permanent) IsTapped() bool {
	return p.tapped
}

func (p Permanent) HasSummoningSickness() bool {
	return p.summoningSickness
}

func (p Permanent) RemoveSummoningSickness() Permanent {
	p.summoningSickness = false
	return p
}

func (p Permanent) Loyalty() int {
	return p.loyalty
}

func (p Permanent) ManaCost() cost.ManaCost {
	return p.manaCost
}

func (p Permanent) ManaValue() int {
	return p.manaCost.Amount().Total()
}

func (per Permanent) Match(predicate query.Predicate) bool {
	return predicate(per)
}

func (p Permanent) Name() string {
	return p.name
}

func (p Permanent) Owner() string {
	return p.owner
}

func (p Permanent) Power() int {
	return p.power
}

func (p Permanent) RulesText() string {
	return p.rulesText
}

func (p Permanent) StaticAbilities() []staticability.StaticAbility {
	return p.staticAbilities
}

func (p Permanent) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range p.staticAbilities {
		if ability.StaticKeyword() != "" {
			keywords = append(keywords, ability.StaticKeyword())
		}
	}
	return keywords
}

func (p Permanent) Subtypes() []mtg.Subtype {
	return p.subtypes
}

func (p Permanent) Supertypes() []mtg.Supertype {
	return p.supertypes
}

// Tap taps the permanent. Returns an error if the permanent is already
// tapped.
// TODO: There are some specifics that matter for tapping tapped creatures,
// but for now we will just disallow it.
func (p Permanent) Tap() (Permanent, error) {
	if p.tapped {
		return p, mtg.ErrAlreadyTapped
	}
	p.tapped = true
	return p, nil
}

func (p Permanent) Toughness() int {
	return p.toughness
}

func (p Permanent) TriggeredAbilities() []Ability {
	return p.triggeredAbilities
}

// Untap untaps the permanent. It is a valid operation even if the permanent
// is already untapped.
// TODO: There are some specifics that matter for untapping untapped
// creatures, but for now we will just allow it.
func (p Permanent) Untap() Permanent {
	p.tapped = false
	return p
}
