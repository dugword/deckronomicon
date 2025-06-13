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
		//activatedAbilities: []core.Ability{},
		// TODO: Maybe parse this into a cost type?

		id:         id,
		controller: playerID,
		owner:      playerID,
		definition: definition,
		loyalty:    definition.Loyalty,
		name:       definition.Name,
		power:      definition.Power,
		rulesText:  definition.RulesText,
		//staticAbilities:    []*StaticAbility{},
		toughness: definition.Toughness,
		// TODO: Do I need to copies these over, or should I save the
		// definition and build them when needed?
		//activatedAbilitySpecs: definition.ActivatedAbilitySpecs,
		//spellAbilitySpec:      definition.SpellAbilitySpec,
		//staticAbilitySpecs:    definition.StaticAbilitySpecs,
		//triggeredAbilitySpecs: definition.TriggeredAbilitySpecs,
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
	/*
		manaCost, err := cost.ParseManaCost(definition.ManaCost)
		if err != nil {
			return Card{}, fmt.Errorf("failed to parse mana cost: %w", err)
		}
	*/
	// card.manaCost = manaCost
	var cardTypes []mtg.CardType
	for _, cardType := range definition.CardTypes {
		cardType, err := mtg.StringToCardType(cardType)
		if err != nil {
			return Card{}, fmt.Errorf("failed to parse card type %q: %w", cardType, err)
		}
		cardTypes = append(cardTypes, cardType)
	}
	card.cardTypes = cardTypes
	var subtypes []mtg.Subtype
	for _, subtype := range definition.Subtypes {
		subtype, ok := mtg.StringToSubtype(subtype)
		if !ok {
			return Card{}, fmt.Errorf("invalid subtype %q", subtype)
		}
		subtypes = append(subtypes, subtype)
	}
	card.subtypes = subtypes
	var supertypes []mtg.Supertype
	for _, supertype := range definition.Supertypes {
		supertype, err := mtg.StringToSupertype(supertype)
		if err != nil {
			return Card{}, fmt.Errorf("failed to parse supertype %q: %w", supertype, err)
		}
		supertypes = append(supertypes, supertype)
	}
	card.supertypes = supertypes
	var ok bool
	var abilityID int
	for _, spec := range definition.ActivatedAbilitySpecs {
		abilityID++
		speed := mtg.SpeedInstant
		if spec.Speed != "" {
			speed, err = mtg.StringToSpeed(spec.Speed)
			if err != nil {
				return Card{}, fmt.Errorf("failed to parse speed %q: %w", spec.Speed, err)
			}
		}
		zone := mtg.ZoneBattlefield
		if spec.Zone != "" {
			zone, ok = mtg.StringToZone(spec.Zone)
			if !ok {
				return Card{}, fmt.Errorf("invalid zone%q", spec.Zone)
			}
		}
		ability := Ability{
			cost: spec.Cost,
			name: spec.Name,
			//effectSpecs: spec.EffectSpecs,
			effects: spec.EffectSpecs,
			id:      fmt.Sprintf("%s-%d", id, abilityID),
			zone:    zone,
			speed:   speed,
			source:  card,
		}
		/*
			var effects []Effect
			// TODO: I think this should happen in the reducer not on load
			for _, effectSpec := range spec.EffectSpecs {
				var modifiers []Tag
				for _, modifier := range effectSpec.Modifiers {
					// TODO: Check if the modifier is valid for the effect
					modifier := Tag{
						Key:   modifier.Key,
						Value: modifier.Value,
					}
					modifiers = append(modifiers, modifier)
				}
				// TODO: Check if required modifiers are present
				var tags []Tag
				// TODO load the tags from the effect spec
				effect := Effect{
					// TODO: Check if the effect exists
					name:      effectSpec.Name,
					modifiers: modifiers,
					optional:  effectSpec.Optional,
					tags:      tags,
				}
				effects = append(effects, effect)
			}
			ability.effects = effects
		*/
		card.activatedAbilities = append(card.activatedAbilities, ability)
	}
	for _, spec := range definition.SpellAbilitySpec.EffectSpecs {
		/*
			var modifiers []Tag
			for _, modifier := range spec.Modifiers {
				modifier := Tag{
					Key:   modifier.Key,
					Value: modifier.Value,
				}
				modifiers = append(modifiers, modifier)
			}
		*/
		card.spellAbility = append(card.spellAbility, spec)
	}

	/*
		for _, spec := range card.activatedAbilitySpecs {
			// TODO: Do this check here or when I am checking for abilities to
			// activate?
			// TODO: Handle the types better, maybe make the spec a type
			if spec.Zone == string(mtg.ZoneHand) || spec.Zone == string(mtg.ZoneGraveyard) {
				// TODO: rename to just Build?
				ability, err := BuildActivatedAbility(state, *spec, &card)
				if err != nil {
					return nil, fmt.Errorf("failed to build activated ability: %w", err)
				}
				card.activatedAbilities = append(card.activatedAbilities, ability)
			}
		}
		for _, spec := range card.staticAbilitySpecs {
			staticAbility, err := BuildStaticAbility(*spec, &card)
			if err != nil {
				return nil, fmt.Errorf("failed to build static ability: %w", err)
			}
			card.staticAbilities = append(card.staticAbilities, staticAbility)
		}
	*/
	return card, nil
}
