package game

import (
	"fmt"
)

// Card represents a card in the game. It contains information about the
// card's attributes which are saved in the underlying object field.
type Card struct {
	object
}

// NewCardFromImport creates a new Card instance from a CardImport.
func NewCardFromImport(cardImport CardImport) (*Card, error) {
	card := Card{
		object: object{
			ActivatedAbilitiesSpec: cardImport.ActivatedAbilitiesSpec,
			CardTypes:              cardImport.CardTypes,
			Loyalty:                cardImport.Loyalty,
			Name:                   cardImport.Name,
			Power:                  cardImport.Power,
			RulesText:              cardImport.RulesText,
			SpellAbilitySpec:       cardImport.SpellAbilitySpec,
			Subtypes:               cardImport.Subtypes,
			Supertypes:             cardImport.Supertypes,
			Toughness:              cardImport.Toughness,
		},
	}
	manaCost, err := ParseManaCost(cardImport.ManaCost)
	if err != nil {
		return nil, fmt.Errorf("failed to create mana cost: %w", err)
	}
	card.object.ManaCost = manaCost
	return &card, nil
}

// CardImport is a struct used for importing card data. It contains fields
// that are used to initialize a Card instance from the imported data.
type CardImport struct {
	ActivatedAbilitiesSpec []ActivatedAbilitySpec `json:"ActivatedAbilities,omitempty"`
	CardTypes              []CardType             `json:"CardTypes,omitempty"`
	Colors                 []string               `json:"Color,omitempty"`
	Loyalty                int                    `json:"Loyalty,omitempty"`
	ManaCost               string                 `json:"ManaCost,omitempty"`
	Name                   string                 `json:"Name,omitempty"`
	Power                  int                    `json:"Power,omitempty"`
	RulesText              string                 `json:"RulesText,omitempty"`
	SpellAbilitySpec       *SpellAbilitySpec      `json:"SpellAbility,omitempty"`
	StaticAbilitiesSpec    []StaticAbilitySpec    `json:"StaticAbilities,omitempty"`
	Subtypes               []Subtype              `json:"Subtypes,omitempty"`
	Supertypes             []Supertype            `json:"Supertypes,omitempty"`
	Toughness              int                    `json:"Toughness,omitempty"`
}

// IsPermanent checks if the card is a permanent. A card is considered a
// permanent if it has one of the following types:
// Artifact
// Battle
// Creature
// Enchantment
// Land
// Planeswalker
func (c *Card) IsPermanent() bool {
	for _, cardType := range c.CardTypes() {
		if cardType.IsPermanent() {
			return true
		}
	}
	return false
}

// IsSpell checks if the card is a spell. A card is considered a spell if it
// has one of the following types:
// Instant
// Sorcery
func (c *Card) IsSpell() bool {
	for _, cardType := range c.CardTypes() {
		if cardType.IsSpell() {
			return true
		}
	}
	return false
}

// ActivatedAbilities returns the activated abilities of the card.
// TODO: This should represent abilities that can be played from hand like
// "Cycling" or from the graveyard like "Flashback".
// TODO: Need to support zone awareness for activated abilities.
func (c *Card) ActivatedAbilities() []ActivatedAbility {
	return c.object.ActivatedAbilities
}

// Card returns the Card instance itself. This is so the Card type implements
// the GameObject interface.
func (c *Card) Card() *Card {
	return c
}

// CardTypes returns the card types of the card. A card can have multiple
// types, such as "Creature" or "Artifact". This method returns a slice of
// CardType representing the types of the card.
func (c *Card) CardTypes() []CardType {
	return c.object.CardTypes
}

// Colors returns the colors of the card. A card can have multiple colors,
// such as "R" or "U". This method returns a slice of string
// representing the colors of the card.
// TODO: Use typed constants for colors instead of strings.
func (c *Card) Colors() *Colors {
	return c.object.Colors
}

// HasType checks if the card has a specific type. It returns true if the
// card has the specified type, and false otherwise.
func (c *Card) HasType(cardType CardType) bool {
	return c.object.HasType(cardType)
}

// ID returns the ID of the card. The ID is a unique identifier for the card
// in the game and remamains the same as the card moves through zones.
func (c *Card) ID() string {
	return c.object.ID
}

// Loyalty returns the loyalty of the card. This is used for planeswalker
// cards. As a card the Loyalty is the starting loyalty of the planeswalker,
// not the current value of the permanent on the battlefield.
func (c *Card) Loyalty() int {
	return c.object.Loyalty
}

// ManaCost returns the mana cost of the card.
func (c *Card) ManaCost() *ManaCost {
	return c.object.ManaCost
}

// Name returns the name of the card.
func (c *Card) Name() string {
	return c.object.Name
}

// Power returns the power of the card. This is used for creature and vehicle
// cards. As a card the Power is the starting power of the creature, not the
// current value of the permanent on the battlefield.
func (c *Card) Power() int {
	return c.object.Power
}

// RulesText returns the rules text of the card. This is used for displaying
// the rules text of the card in the game. The rules text is a string that
// describes the abilities and effects of the card. It is not used for any in
// game logic.
func (c *Card) RulesText() string {
	return c.object.RulesText
}

// StaticAbilities returns the static abilities of the card. NOTE: this is not
// the same as the static abilities of a permanent. A card can have static
// abilities that are not present on the permanent. For example, a card can
// say "You may cast this card as though it had flash', which is a static
// ability enabling the card to be cast at instant speed.
func (c *Card) StaticAbilities() []StaticAbility {
	return c.object.StaticAbilities
}

// Subtypes returns the subtypes of the card.
func (c *Card) Subtypes() []Subtype {
	return c.object.Subtypes
}

// Supertypes returns the supertypes of the card.
func (c *Card) Supertypes() []Supertype {
	return c.object.Supertypes
}

// Toughness returns the toughness of the card. This is used for creature and
// vehicle cards. As a card the Toughness is the starting toughness of the
// creatue, not the current value of the permanent on the battlefield.
func (c *Card) Toughness() int {
	return c.object.Toughness
}
