package gob

import (

	// "deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"fmt"
)

// Permanent represents a permanent card on the battlefield.
type Permanent struct {
	activatedAbilities []Ability
	card               Card
	cardTypes          []mtg.CardType
	colors             mtg.Colors
	controller         string
	id                 string
	loyalty            int
	manaCost           string
	name               string
	owner              string
	power              int
	rulesText          string
	staticAbilities    []StaticAbility
	subtypes           []mtg.Subtype
	summoningSickness  bool
	supertypes         []mtg.Supertype
	tapped             bool
	toughness          int
	triggeredAbilities []Ability
}

// NewPermanent creates a new Permanent instance from a Card.
func NewPermanent(id string, card Card, playerID string) (Permanent, error) {
	permanent := Permanent{
		// TODO: Do I need to generate new IDs for the abilities?
		card:       card,
		cardTypes:  card.CardTypes(),
		colors:     card.Colors(),
		controller: playerID,
		id:         id,
		loyalty:    card.Loyalty(),
		//manaCost:           card.ManaCost(),
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
	if card.Match(has.CardType(mtg.CardTypeCreature)) {
		permanent.summoningSickness = true
	}
	var abilityID int
	fmt.Println("Building permanent:", permanent.name, "ID:", id)
	for _, a := range card.ActivatedAbilities() {
		fmt.Println("Building activated ability:", a.Name())
		abilityID++
		ability := Ability{
			effects: a.Effects(),
			name:    a.Name(),
			cost:    a.Cost(),
			id:      fmt.Sprintf("%s-%d", id, abilityID),
			zone:    a.zone,
			speed:   a.speed,
			source:  permanent,
		}
		for _, effect := range a.Effects() {
			fmt.Println("Adding effect to ability:", effect.Name())
		}
		fmt.Printf("Ability => %+v\n", ability)
		permanent.activatedAbilities = append(permanent.activatedAbilities, ability)
	}
	/*
		for _, spec := range card.ActivatedAbilitySpecs() {
			// TODO: use better types
			if spec.Zone == string(mtg.ZoneHand) || spec.Zone == string(mtg.ZoneGraveyard) {
				continue
			}
			ability, err := BuildActivatedAbility(state, *spec, &permanent)
			if err != nil {
				return nil, fmt.Errorf("failed to build activated ability: %w", err)
			}
			permanent.activatedAbilities = append(permanent.activatedAbilities, ability)
		}
	*/
	return permanent, nil
}

// ActivatedAbilities returns the activated abilities of the permanent.
func (p Permanent) ActivatedAbilities() []Ability {
	return p.activatedAbilities
}

// Card returns the card associated with the permanent.
func (p Permanent) Card() Card {
	return p.card
}

// CardTypes returns the card types of the permanent.
func (p Permanent) CardTypes() []mtg.CardType {
	return p.cardTypes
}

// Colors returns the colors of the permanent.
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

// IsTapped checks if the permanent is tapped.
func (p Permanent) IsTapped() bool {
	return p.tapped
}

// HasSummoningSickness checks if the permanent has summoning sickness.
func (p Permanent) HasSummoningSickness() bool {
	return p.summoningSickness
}

// Loyalty returns the loyalty of the permanent.
func (p Permanent) Loyalty() int {
	return p.loyalty
}

/*
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
*/

func (per Permanent) Match(p query.Predicate) bool {
	return p(per)
}

// Name returns the name of the permanent.
func (p Permanent) Name() string {
	return p.name
}

// Power returns the power of the permanent.
func (p Permanent) Power() int {
	return p.power
}

// RemoveSummoningSickness removes summoning sickness from the permanent. It
// is a valid operation even if the permanent does not have summoning
// sickness.
/*
func (p Permanent) RemoveSummoningSickness() {
	p.summoningSickness = false
}
*/

// RulesText returns the rules text of the permanent. The RulesText does not
// impact the game logic.
func (p Permanent) RulesText() string {
	return p.rulesText
}

// StaticAbilities returns the static abilities of the permanent.
func (p Permanent) StaticAbilities() []StaticAbility {
	return p.staticAbilities
}

// StaticAbilities returns the static abilities of the permanent.
func (p Permanent) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range p.staticAbilities {
		keyword, ok := mtg.StringToStaticKeyword(ability.ID())
		if !ok {
			continue
		}
		keywords = append(keywords, keyword)
	}
	return keywords
}

// Subtypes returns the subtypes of the permanent.
func (p Permanent) Subtypes() []mtg.Subtype {
	return p.subtypes
}

// Supertypes returns the supertypes of the permanent.
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

// Toughness returns the toughness of the permanent.
func (p Permanent) Toughness() int {
	return p.toughness
}

// TriggeredAbilities returns the triggered abilities of the permanent.
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
