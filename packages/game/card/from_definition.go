package card

import (
	"deckronomicon/packages/game/ability/activated"
	"deckronomicon/packages/game/ability/static"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

// NewCardFromCardData creates a new Card instance from the given CardData.
// It initializes the card's attributes, including its abilities, colors,
// mana cost, types, and supertypes. It also generates a unique ID for the
// card. The function returns a pointer to the newly created Card instance
// and an error if any occurred during the creation process. Any Activated
// Abilities or Static Abilities that are available while the card is in the
// players hand or graveyard are built and added to the card.
func NewCardFromCardDefinition(state state, definition definition.Card) (*Card, error) {
	card := Card{
		activatedAbilities: []*activated.Ability{},
		definition:         &definition,
		loyalty:            definition.Loyalty,
		name:               definition.Name,
		power:              definition.Power,
		rulesText:          definition.RulesText,
		staticAbilities:    []*static.Ability{},
		toughness:          definition.Toughness,
		// TODO: Do I need to copies these over, or should I save the
		// definition and build them when needed?
		activatedAbilitySpecs: definition.ActivatedAbilitySpecs,
		spellAbilitySpec:      definition.SpellAbilitySpec,
		triggeredAbilitySpecs: definition.TriggeredAbilitySpecs,
	}
	cardColors, err := mtg.StringsToColors(definition.Colors)
	if err != nil {
		return nil, fmt.Errorf("failed to create colors: %w", err)
	}
	card.colors = cardColors
	manaCost, err := cost.ParseManaCost(definition.ManaCost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana cost: %w", err)
	}
	card.manaCost = manaCost
	var cardTypes []mtg.CardType
	for _, cardType := range definition.CardTypes {
		cardType, err := mtg.StringToCardType(cardType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse card types: %w", err)
		}
		cardTypes = append(cardTypes, cardType)
	}
	card.cardTypes = cardTypes
	var subtypes []mtg.Subtype
	for _, subtype := range definition.Subtypes {
		subtype, err := mtg.StringToSubtype(subtype)
		if err != nil {
			return nil, fmt.Errorf("failed to parse subtypes: %w", err)
		}
		subtypes = append(subtypes, subtype)
	}
	card.subtypes = subtypes
	var supertypes []mtg.Supertype
	for _, supertype := range definition.Supertypes {
		supertype, err := mtg.StringToSupertype(supertype)
		if err != nil {
			return nil, fmt.Errorf("failed to parse supertypes: %w", err)
		}
		supertypes = append(supertypes, supertype)
	}
	card.supertypes = supertypes
	card.id = state.GetNextID()
	return &card, nil
}
