package game

import "fmt"

// Card represents a card in the game. It contains all the information about
// the card, including its name, mana cost, power, toughness, and abilities.
// The Card type is used to represent a card in the game and holds the specs
// to create additional game state. The Card type is also used to represent a
// card in the player's hand, graveyard, and library. The Card type is used to
// create Spells on the stack and Permanents on the battlefield. The
// underlying card data is preserved in the CardData struct, if card
// attributes are ever overwritten via stickers or other effects. Each card in
// a player's library is a unique instance of the Card type, even if the card
// is a copy of another card.
type Card struct {
	activatedAbilities    []*ActivatedAbility
	activatedAbilitySpecs []*ActivatedAbilitySpec
	cardData              *CardData
	cardTypes             []CardType
	colors                *Colors
	id                    string
	loyalty               int
	manaCost              *ManaCost
	name                  string
	power                 int
	rulesText             string
	spellAbilitySpec      *SpellAbilitySpec
	staticAbilities       []*StaticAbility
	staticAbilitySpecs    []*StaticAbilitySpec
	triggeredAbilitySpecs []*TriggeredAbilitySpec
	subtypes              []Subtype
	supertypes            []Supertype
	toughness             int
}

// NewCardFromCardData creates a new Card instance from the given CardData.
// It initializes the card's attributes, including its abilities, colors,
// mana cost, types, and supertypes. It also generates a unique ID for the
// card. The function returns a pointer to the newly created Card instance
// and an error if any occurred during the creation process. Any Activated
// Abilities or Static Abilities that are available while the card is in the
// players hand or graveyard are built and added to the card.
func NewCardFromCardData(cardData CardData) (*Card, error) {
	card := Card{
		activatedAbilities:    []*ActivatedAbility{},
		activatedAbilitySpecs: cardData.ActivatedAbilitySpecs,
		cardData:              &cardData,
		loyalty:               cardData.Loyalty,
		name:                  cardData.Name,
		power:                 cardData.Power,
		rulesText:             cardData.RulesText,
		spellAbilitySpec:      cardData.SpellAbilitySpec,
		staticAbilitySpecs:    cardData.StaticAbilitySpecs,
		staticAbilities:       []*StaticAbility{},
		toughness:             cardData.Toughness,
		triggeredAbilitySpecs: cardData.TriggeredAbilitySpecs,
	}
	colors, err := StringsToColors(cardData.Colors)
	if err != nil {
		return nil, fmt.Errorf("failed to create colors: %w", err)
	}
	card.colors = colors
	manaCost, err := ParseManaCost(cardData.ManaCost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana cost: %w", err)
	}
	card.manaCost = manaCost
	var cardTypes []CardType
	for _, cardType := range cardData.CardTypes {
		cardType, err := StringToCardType(cardType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse card types: %w", err)
		}
		cardTypes = append(cardTypes, cardType)
	}
	card.cardTypes = cardTypes
	var subtypes []Subtype
	for _, subtype := range cardData.Subtypes {
		subtype, err := StringToSubtype(subtype)
		if err != nil {
			return nil, fmt.Errorf("failed to parse subtypes: %w", err)
		}
		subtypes = append(subtypes, subtype)
	}
	card.subtypes = subtypes
	var supertypes []Supertype
	for _, supertype := range cardData.Supertypes {
		supertype, err := StringToSupertype(supertype)
		if err != nil {
			return nil, fmt.Errorf("failed to parse supertypes: %w", err)
		}
		supertypes = append(supertypes, supertype)
	}
	card.supertypes = supertypes
	card.id = GetNextID()
	for _, spec := range card.activatedAbilitySpecs {
		if spec.Zone == ZoneHand || spec.Zone == ZoneGraveyard {
			ability, err := BuildActivatedAbility(*spec, &card)
			if err != nil {
				return nil, fmt.Errorf("failed to build activated ability: %w", err)
			}
			card.activatedAbilities = append(card.activatedAbilities, ability)
		}
	}
	for _, spec := range card.staticAbilitySpecs {
		if spec.Zone == ZoneHand || spec.Zone == ZoneGraveyard {
			ability, err := BuildStaticAbility(*spec, &card)
			if err != nil {
				return nil, fmt.Errorf("failed to build activated ability: %w", err)
			}
			card.staticAbilities = append(card.staticAbilities, ability)
		}
	}
	return &card, nil
}

// IsPermanent checks if the card is a permanent. A card is considered a
// permanent if it has one of the following types:
//   - Artifact
//   - Battle
//   - Creature
//   - Enchantment
//   - Land
//   - Planeswalker
func (c *Card) IsPermanent() bool {
	for _, cardType := range c.CardTypes() {
		if cardType.IsPermanent() {
			return true
		}
	}
	return false
}

// IsLand checks if the card is a land. A card is considered a land if it
// has the Land CardType
func (c *Card) IsLand() bool {
	for _, cardType := range c.CardTypes() {
		if cardType.IsLand() {
			return true
		}
	}
	return false
}

// IsSpell checks if the card is a spell. A card is considered a spell if it
// has one of the following types:
//   - Instant
//   - Sorcery
func (c *Card) IsSpell() bool {
	for _, cardType := range c.CardTypes() {
		if cardType.IsSpell() {
			return true
		}
	}
	return false
}

// ActivatedAbilities returns the activated abilities of the card available.
// NOTE: These are the activated abilities of the card itself, not the
// activated abilities of the permanent. A card can have activated abilities
// that are not present on the permanent.
func (c *Card) ActivatedAbilities() []*ActivatedAbility {
	return c.activatedAbilities
}

// CardTypes returns the card types of the card. A card can have multiple
// types, such as "Creature" or "Artifact". This method returns a slice of
// CardType representing the types of the card.
func (c *Card) CardTypes() []CardType {
	return c.cardTypes
}

// Colors returns the colors of the card. A card can have multiple colors. This
// method returns a Colors struct representing all colors the card is.
func (c *Card) Colors() *Colors {
	return c.colors
}

// HasCardType checks if the card has a specific type. It returns true if the
// card has the specified type, and false otherwise.
func (c *Card) HasCardType(cardType CardType) bool {
	for _, t := range c.cardTypes {
		if t == cardType {
			return true
		}
	}
	return false
}

// HasSubtype checks if the card has a specific type. It returns true if the
// card has the specified type, and false otherwise.
func (c *Card) HasSubtype(subtype Subtype) bool {
	for _, t := range c.subtypes {
		if t == subtype {
			return true
		}
	}
	return false
}

// ID returns the ID of the card. The ID is a unique identifier for the card
// in the game and remamains the same as the card moves through zones. The ID
// is unique per card instance, so two copies of the same card will have
// different IDs.
func (c *Card) ID() string {
	return c.id
}

// Loyalty returns the loyalty of the card. This is used for planeswalker
// cards. As a card the Loyalty is the starting loyalty of the planeswalker,
// not the current value of the permanent on the battlefield.
func (c *Card) Loyalty() int {
	return c.loyalty
}

// ManaCost returns the mana cost of the card.
func (c *Card) ManaCost() *ManaCost {
	return c.manaCost
}

// Name returns the name of the card.
func (c *Card) Name() string {
	return c.name
}

// Power returns the power of the card. This is used for creature and vehicle
// cards. As a card the Power is the starting power of the creature, not the
// current value of the permanent on the battlefield.
func (c *Card) Power() int {
	return c.power
}

// RulesText returns the rules text of the card. This is used for displaying
// the rules text of the card in the game. The rules text is a string that
// describes the abilities and effects of the card. It is not used for any in
// game logic.
func (c *Card) RulesText() string {
	return c.rulesText
}

// StaticAbilities returns the static abilities of the card available in the
// provide zone
// NOTE: this is not the same as the static abilities of a permanent. A card
// can have static abilities that are not present on the permanent.
func (c *Card) StaticAbilities(zone string) []*StaticAbility {
	var abilities []*StaticAbility
	for _, ability := range c.staticAbilities {
		if ability.Zone == zone {
			abilities = append(abilities, ability)
		}
	}
	return abilities
}

// Subtypes returns the subtypes of the card.
func (c *Card) Subtypes() []Subtype {
	return c.subtypes
}

// Supertypes returns the supertypes of the card.
func (c *Card) Supertypes() []Supertype {
	return c.supertypes
}

// Toughness returns the toughness of the card. This is used for creature and
// vehicle cards. As a card the Toughness is the starting toughness of the
// creatue, not the current value of the permanent on the battlefield.
func (c *Card) Toughness() int {
	return c.toughness
}
