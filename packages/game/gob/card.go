package gob

import (
	// "deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

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
	activatedAbilities []Ability
	// activatedAbilitySpecs []definition.ActivatedAbilitySpec
	definition   definition.Card
	cardTypes    []mtg.CardType
	colors       mtg.Colors
	id           string
	loyalty      int
	manaCost     string
	name         string
	power        int
	rulesText    string
	spellAbility []Effect
	// TODO: Maybe this should just be spell spec? sepll effect spec?
	//spellAbilitySpec      definition.SpellAbilitySpec
	staticAbilities []StaticAbility
	//staticAbilitySpecs    []definition.StaticAbilitySpec
	//triggeredAbilitySpecs []definition.TriggeredAbilitySpec
	subtypes   []mtg.Subtype
	supertypes []mtg.Supertype
	toughness  int
}

func NewCard(id, name string) Card {
	card := Card{
		id:   id,
		name: name,
	}
	return card
}

// ActivatedAbilities returns the activated abilities of the card available.
// NOTE: These are the activated abilities of the card itself, not the
// activated abilities of the permanent. A card can have activated abilities
// that are not present on the permanent.
func (c Card) ActivatedAbilities() []Ability {
	return c.activatedAbilities
}

/*
func (c Card) ActivatedAbilitySpecs() []definition.ActivatedAbilitySpec {
	return c.activatedAbilitySpecs
}
*/

// CardTypes returns the card types of the card. A card can have multiple
// types, such as "Creature" or "Artifact". This method returns a slice of
// CardType representing the types of the card.
func (c Card) CardTypes() []mtg.CardType {
	return c.cardTypes
}

// Colors returns the colors of the card. A card can have multiple colors. This
// method returns a Colors struct representing all colors the card is.
func (c Card) Colors() mtg.Colors {
	return c.colors
}

// TODO: Is this the best way to do this?
func (c Card) Description() string {
	return c.rulesText
}

// ID returns the ID of the card. The ID is a unique identifier for the card
// in the game and remamains the same as the card moves through zones. The ID
// is unique per card instance, so two copies of the same card will have
// different IDs.
func (c Card) ID() string {
	return c.id
}

func (c Card) Match(predicate query.Predicate) bool {
	return predicate(c)
}

// Loyalty returns the loyalty of the card. This is used for planeswalker
// cards. As a card the Loyalty is the starting loyalty of the planeswalker,
// not the current value of the permanent on the battlefield.
func (c Card) Loyalty() int {
	return c.loyalty
}

// ManaCost returns the mana cost of the card.
func (c Card) ManaCost() string {
	return c.manaCost
}

/*
func (c Card) ManaValue() int {
	return c.ManaCost().ManaValue()
}
*/

// Name returns the name of the card.
func (c Card) Name() string {
	return c.name
}

// Power returns the power of the card. This is used for creature and vehicle
// cards. As a card the Power is the starting power of the creature, not the
// current value of the permanent on the battlefield.
func (c Card) Power() int {
	return c.power
}

// RulesText returns the rules text of the card. This is used for displaying
// the rules text of the card in the game. The rules text is a string that
// describes the abilities and effects of the card. It is not used for any in
// game logic.
func (c Card) RulesText() string {
	return c.rulesText
}

func (c Card) SpellAbility() []Effect {
	return c.spellAbility
}

// StaticAbilities returns the static abilities of the card.
func (c Card) StaticAbilities() []StaticAbility {
	var abilities = append([]StaticAbility{}, c.staticAbilities...)
	return abilities
}

func (c Card) StaticKeywords() []mtg.StaticKeyword {
	var keywords []mtg.StaticKeyword
	for _, ability := range c.staticAbilities {
		/*
			if ability == nil {
				continue
			}
		*/
		// TODO: Check this error or just use typed values all the way down.
		keyword, _ := mtg.StringToStaticKeyword(ability.ID())
		keywords = append(keywords, keyword)
	}
	return keywords
}

// Subtypes returns the subtypes of the card.
func (c Card) Subtypes() []mtg.Subtype {
	return c.subtypes
}

// Supertypes returns the supertypes of the card.
func (c Card) Supertypes() []mtg.Supertype {
	return c.supertypes
}

// Toughness returns the toughness of the card. This is used for creature and
// vehicle cards. As a card the Toughness is the starting toughness of the
// creatue, not the current value of the permanent on the battlefield.
func (c Card) Toughness() int {
	return c.toughness
}
