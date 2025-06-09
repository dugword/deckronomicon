package gob

import (
	//"deckronomicon/packages/game/core"
	// "deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"

	// "deckronomicon/packages_old/game/core"
	"fmt"
	"strings"
)

// TODO make these interfaces private
/*
type BattlefieldAdder interface {
	AddToBattlefield(permanent Permanent)
}

type CardPlacer interface {
	PlaceCard(card Card, zone mtg.Zone) error
}
*/

// Spell represents a spell object on the stack.
type Spell struct {
	card      Card
	cardTypes []mtg.CardType
	colors    mtg.Colors
	//flashback bool
	id      string
	loyalty int
	//manaCost        *cost.ManaCost
	name            string
	power           int
	rulesText       string
	effects         []Effect
	staticAbilities []StaticAbility
	subtypes        []mtg.Subtype
	supertypes      []mtg.Supertype
	toughness       int
}

// NewSpell creates a new Spell instance from a Card.
func NewSpell(id string, card Card) (Spell, error) {
	spell := Spell{
		card:      card,
		cardTypes: card.CardTypes(),
		colors:    card.Colors(),
		id:        id,
		loyalty:   card.Loyalty(),
		//manaCost:        card.ManaCost(),
		name:            card.Name(),
		power:           card.Power(),
		rulesText:       card.RulesText(),
		staticAbilities: card.StaticAbilities(),
		subtypes:        card.Subtypes(),
		supertypes:      card.Supertypes(),
		toughness:       card.Toughness(),
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

// Effects returns the effects of the spell.
func (s Spell) Effects() []Effect {
	return s.effects
}

// Description returns a string representation of the activated ability.
func (s Spell) Description() string {
	var descriptions []string
	for _, effect := range s.effects {

		// TODO: Come up with a better way to handle descriptions
		descriptions = append(descriptions, effect.Name())
		// descriptions = append(descriptions, effect.Description())

	}
	// TODO: Support additional costs
	// return fmt.Sprintf("%s: %s", s.ManaCost().Description(), strings.Join(descriptions, ", "))
	return fmt.Sprintf("%s: %s", "<cost>", strings.Join(descriptions, ", "))
}

// ID returns the ID of the spell.
func (s Spell) ID() string {
	return s.id
}

// Loyalty returns the loyalty of the spell.
func (s Spell) Loyalty() int {
	return s.loyalty
}

/*
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
*/

func (s Spell) Match(p query.Predicate) bool {
	return p(s)
}

// Name returns the name of the spell.
func (s Spell) Name() string {
	return s.name
}

// Power returns the power of the spell.
func (s Spell) Power() int {
	return s.power
}

/*
func (s *Spell) Resolve(state ore.State, player core.Player) error {
	// TODO: Do I need to check for effects here?
	if s.Effects == nil && s.card.Match(is.Permanent()) {
		permanent, err := NewPermanent(s.card, state, player)
		if err != nil {
			return fmt.Errorf("failed to create permanent: %w", err)
		}
		battlefieldAdder, ok := state.(BattlefieldAdder)
		if !ok {
			return fmt.Errorf("state does not implement BattlefieldAdder")
		}
		battlefieldAdder.AddToBattlefield(permanent)
	}
	for _, effect := range s.effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	cardPlacer, ok := state.(CardPlacer)
	if !ok {
		return fmt.Errorf("state does not implement CardPlacer")
	}
	if s.flashback {
		if err := cardPlacer.PlaceCard(s.card, mtg.ZoneExile); err != nil {
			return fmt.Errorf("cannot place card in exile: %w", err)
		}
		return nil
	}
	if err := cardPlacer.PlaceCard(s.card, mtg.ZoneGraveyard); err != nil {
		return fmt.Errorf("cannot move spell to graveyard: %w", err)
	}
	return nil
}
*/

// RulesText returns the rules text of the spell. The RulesText does not
// impact the game logic.
func (s Spell) RulesText() string {
	return s.rulesText
}

// StaticAbilities returns the static abilities of the spell
func (s Spell) StaticAbilities() []StaticAbility {
	return s.staticAbilities
}

func (s Spell) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range s.staticAbilities {
		keyword, ok := mtg.StringToStaticKeyword(ability.ID())
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

/*
func (s *Spell) Splice(state core.State, card *Card) error {
	for _, effect := range card.SpellAbility() {
		s.effects = append(s.effects, effect)
	}
	return nil
}
*/
