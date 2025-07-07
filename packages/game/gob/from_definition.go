package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"fmt"
)

func NewCardFromDefinition(definition *definition.Card) *Card {
	card := Card{
		id:         definition.ID,
		controller: definition.Controller,
		owner:      definition.Owner,
		loyalty:    definition.Loyalty,
		name:       definition.Name,
		power:      definition.Power,
		rulesText:  definition.RulesText,
		toughness:  definition.Toughness,
	}
	manaCost, err := cost.ParseMana(definition.ManaCost)
	if err != nil {
		panic(fmt.Errorf("failed to parse cost %q: %w", definition.ManaCost, err))
	}
	additionalCost, err := cost.Parse(definition.AdditionalCost)
	if err != nil {
		panic(fmt.Errorf("failed to parse additional cost %q: %w", definition.AdditionalCost, err))
	}
	cardColors, err := mtg.StringsToColors(definition.Colors)
	if err != nil {
		panic(fmt.Errorf("failed to parse colors %s: %w", definition.Colors, err))
	}
	var cardTypes []mtg.CardType
	for _, cardTypeDefinition := range definition.CardTypes {
		cardType, ok := mtg.StringToCardType(cardTypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse card type %q: %w", cardType, err))
		}
		cardTypes = append(cardTypes, cardType)
	}
	var subtypes []mtg.Subtype
	for _, subtypeDefinition := range definition.Subtypes {
		subtype, ok := mtg.StringToSubtype(subtypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse subtype %q: %w", subtypeDefinition, err))
		}
		subtypes = append(subtypes, subtype)
	}
	var supertypes []mtg.Supertype
	for _, supertypeDefinition := range definition.Supertypes {
		supertype, ok := mtg.StringToSupertype(supertypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse supertype %q: %w", supertypeDefinition, err))
		}
		supertypes = append(supertypes, supertype)
	}
	var activatedAbilities []*Ability
	for i, abilityDefinition := range definition.ActivatedAbilities {
		ability, err := NewAbility(&card, i, abilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build activated ability %s: %w", abilityDefinition.Name, err))
		}
		activatedAbilities = append(card.activatedAbilities, ability)
	}
	var spellAbility []effect.Effect
	for _, effectDefinition := range definition.SpellAbility {
		effect, err := effect.New(effectDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build spell effect %s: %w", effectDefinition.Name, err))
		}
		spellAbility = append(spellAbility, effect)
	}
	var staticAbilities []staticability.StaticAbility
	for _, staticAbilityDefinition := range definition.StaticAbilities {
		staticAbility, err := staticability.New(staticAbilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build static ability %s: %w", staticAbilityDefinition.Name, err))
		}
		staticAbilities = append(staticAbilities, staticAbility)
	}
	var triggeredAbilities []*TriggeredAbility
	for _, triggeredAbilityDefinition := range definition.TriggeredAbilities {
		triggeredAbility, err := NewTriggeredAbility(&card, len(triggeredAbilities), triggeredAbilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build triggered ability %s: %w", triggeredAbilityDefinition.Name, err))
		}
		triggeredAbilities = append(triggeredAbilities, triggeredAbility)
	}
	card.activatedAbilities = activatedAbilities
	card.cardTypes = cardTypes
	card.colors = cardColors
	card.manaCost = manaCost
	card.additionalCost = additionalCost
	card.spellAbility = spellAbility
	card.staticAbilities = staticAbilities
	card.triggeredAbilities = triggeredAbilities
	card.subtypes = subtypes
	card.supertypes = supertypes
	return &card
}

func NewPermanentFromCardDefinition(
	definition *definition.Card,
) *Permanent {
	card := NewCardFromDefinition(definition)
	permanent, err := NewPermanent(
		fmt.Sprintf("Permanent from %s", definition.ID),
		card,
		definition.Controller,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create permanent from card: %w", err))
	}
	return permanent
}

func NewPermanentFromDefinition(
	definition *definition.Permanent,
) *Permanent {
	permanent := Permanent{
		id:                definition.ID,
		controller:        definition.Controller,
		owner:             definition.Owner,
		loyalty:           definition.Loyalty,
		name:              definition.Name,
		power:             definition.Power,
		rulesText:         definition.RulesText,
		summoningSickness: definition.SummoningSickness,
		tapped:            definition.Tapped,
		toughness:         definition.Toughness,
	}
	if definition.Card != nil {
		permanent.card = NewCardFromDefinition(definition.Card)
	}
	manaCost, err := cost.ParseMana(definition.ManaCost)
	if err != nil {
		panic(fmt.Errorf("failed to parse mana cost %q: %w", definition.ManaCost, err))
	}
	cardColors, err := mtg.StringsToColors(definition.Colors)
	if err != nil {
		panic(fmt.Errorf("failed to parse colors %s: %w", definition.Colors, err))
	}
	var cardTypes []mtg.CardType
	for _, cardTypeDefinition := range definition.CardTypes {
		cardType, ok := mtg.StringToCardType(cardTypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse card type %q: %w", cardType, err))
		}
		cardTypes = append(cardTypes, cardType)
	}
	var subtypes []mtg.Subtype
	for _, subtypeDefinition := range definition.Subtypes {
		subtype, ok := mtg.StringToSubtype(subtypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse subtype %q: %w", subtypeDefinition, err))
		}
		subtypes = append(subtypes, subtype)
	}
	var supertypes []mtg.Supertype
	for _, supertypeDefinition := range definition.Supertypes {
		supertype, ok := mtg.StringToSupertype(supertypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse supertype %q: %w", supertypeDefinition, err))
		}
		supertypes = append(supertypes, supertype)
	}
	var activatedAbilities []*Ability
	for i, abilityDefinition := range definition.ActivatedAbilities {
		ability, err := NewAbility(&permanent, i, abilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build activated ability %s: %w", abilityDefinition.Name, err))
		}
		activatedAbilities = append(activatedAbilities, ability)
	}
	var staticAbilities []staticability.StaticAbility
	for _, staticAbilityDefinition := range definition.StaticAbilities {
		staticAbility, err := staticability.New(staticAbilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build static ability %s: %w", staticAbilityDefinition.Name, err))
		}
		staticAbilities = append(staticAbilities, staticAbility)
	}
	permanent.activatedAbilities = activatedAbilities
	permanent.cardTypes = cardTypes
	permanent.colors = cardColors
	permanent.manaCost = manaCost
	permanent.staticAbilities = staticAbilities
	permanent.subtypes = subtypes
	permanent.supertypes = supertypes
	return &permanent
}

func NewSpellFromCardDefinition(
	spellID,
	playerID string,
	cardDefinition *definition.Card,
	effectWithTargets []*effect.EffectWithTarget,
	flashback bool,
) *Spell {
	card := NewCardFromDefinition(cardDefinition)
	spell, err := NewSpell(spellID, card, playerID, effectWithTargets, flashback)
	if err != nil {
		panic(fmt.Errorf("failed to create spell from card: %w", err))
	}
	return spell
}

func NewSpellFromDefinition(
	definition definition.Spell,
) *Spell {
	spell := Spell{
		id:         definition.ID,
		controller: definition.Controller,
		owner:      definition.Owner,
		flashback:  definition.Flashback,
	}
	if definition.Card != nil {
		spell.card = NewCardFromDefinition(definition.Card)
	}
	manaCost, err := cost.ParseMana(definition.ManaCost)
	if err != nil {
		panic(fmt.Errorf("failed to parse cost %q: %w", definition.ManaCost, err))
	}
	cardColors, err := mtg.StringsToColors(definition.Colors)
	if err != nil {
		panic(fmt.Errorf("failed to parse colors %s: %w", definition.Colors, err))
	}
	var cardTypes []mtg.CardType
	for _, cardTypeDefinition := range definition.CardTypes {
		cardType, ok := mtg.StringToCardType(cardTypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse card type %q: %w", cardType, err))
		}
		cardTypes = append(cardTypes, cardType)
	}
	var subtypes []mtg.Subtype
	for _, subtypeDefinition := range definition.Subtypes {
		subtype, ok := mtg.StringToSubtype(subtypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse subtype %q: %w", subtypeDefinition, err))
		}
		subtypes = append(subtypes, subtype)
	}
	var supertypes []mtg.Supertype
	for _, supertypeDefinition := range definition.Supertypes {
		supertype, ok := mtg.StringToSupertype(supertypeDefinition)
		if !ok {
			panic(fmt.Errorf("failed to parse supertype %q: %w", supertypeDefinition, err))
		}
		supertypes = append(supertypes, supertype)
	}
	var staticAbilities []staticability.StaticAbility
	for _, staticAbilityDefinition := range definition.StaticAbilities {
		staticAbility, err := staticability.New(staticAbilityDefinition)
		if err != nil {
			panic(fmt.Errorf("failed to build static ability %s: %w", staticAbilityDefinition.Name, err))
		}
		staticAbilities = append(staticAbilities, staticAbility)
	}
	spell.cardTypes = cardTypes
	spell.colors = cardColors
	spell.manaCost = manaCost
	spell.staticAbilities = staticAbilities
	spell.subtypes = subtypes
	spell.supertypes = supertypes
	return &spell
}
