package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

// This reads the definition that comes from JSON files and validates all the types and stuff actuatlly
// exist, casts them to the correct types, and then builds the card.

// "deckronomicon/packages/game/core"

// NewCardFromCardData creates a new Card instance from the given CardData.
// It initializes the card's attributes, including its abilities, colors,
// mana cost, types, and supertypes. It also generates a unique ID for the
// card. The function returns a pointer to the newly created Card instance
// and an error if any occurred during the creation process. Any Activated
// Abilities or Static Abilities that are available while the card is in the
// players hand or graveyard are built and added to the card.

// TODO: Make this just NewCard

func NewCardFromCardDefinition(id, playerID string, definition definition.Card) (Card, error) {
	card := Card{
		id:         id,
		controller: playerID,
		owner:      playerID,
		definition: definition,
		loyalty:    definition.Loyalty,
		name:       definition.Name,
		power:      definition.Power,
		rulesText:  definition.RulesText,
		toughness:  definition.Toughness,
		cardTypes:  definition.CardTypes,
		subtypes:   definition.Subtypes,
		supertypes: definition.Supertypes,
	}
	manaCost, err := cost.ParseManaCost(definition.ManaCost)
	if err != nil {
		return Card{}, fmt.Errorf("failed to parse mana cost %q: %w", definition.ManaCost, err)
	}
	card.manaCost = manaCost
	cardColors, err := mtg.StringsToColors(definition.Colors)
	if err != nil {
		return Card{}, fmt.Errorf("failed to parse colors %s: %w", definition.Colors, err)
	}
	card.colors = cardColors

	for i, spec := range definition.ActivatedAbilitySpecs {
		speed := mtg.SpeedInstant
		if spec.Speed != "" {
			speed = spec.Speed
		}
		zone := mtg.ZoneBattlefield
		if spec.Zone != "" {
			zone = spec.Zone
		}
		abilityCost, err := cost.ParseCost(spec.Cost)
		if err != nil {
			return Card{}, fmt.Errorf("failed to parse cost %q: %w", spec.Cost, err)
		}
		ability := Ability{
			cost:        abilityCost,
			name:        spec.Name,
			effectSpecs: spec.EffectSpecs,
			id:          fmt.Sprintf("%s-%d", id, i+1),
			zone:        zone,
			speed:       speed,
			source:      card,
		}
		card.activatedAbilities = append(card.activatedAbilities, ability)
	}
	card.spellAbility = append(card.spellAbility, definition.SpellAbilitySpec...)
	for _, spec := range definition.StaticAbilitySpecs {
		staticCost, err := cost.ParseCost(spec.Cost)
		if err != nil {
			return Card{}, fmt.Errorf("failed to parse static ability cost %q: %w", spec.Cost, err)
		}
		card.staticAbilities = append(card.staticAbilities, StaticAbility{
			name:      spec.Name,
			cost:      staticCost,
			modifiers: spec.Modifiers,
		})
	}
	return card, nil
}
