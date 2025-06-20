package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob/gobtest"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

// These configs enable game objects to be loaded from a structured format. This structure is useful for saving and loading game states
// or for initializing a game with predefined settings.
// It includes all necessary fields to reconstruct a game instance, such as player states, battlefield, stack, and game phase.
// These configurations should not be used for managing game state during active gameplay, but rather for testing, setup or serialization
// purposes

func LoadPermanentFromConfig(config gobtest.PermanentConfig) Permanent {
	permanent := Permanent{
		id:         config.ID,
		controller: config.Controller,
		owner:      config.Owner,
		name:       config.Name,
		card:       LoadCardFromConfig(config.Owner, config.Card),
		tapped:     config.Tapped,
	}
	// TODO: this is redundant with code elsewhere, consider refactoring
	// and making a NewAbility function.
	var activatedAbilities []Ability
	for i, ability := range config.ActivatedAbilities {
		abilityCost, err := cost.ParseManaCost(ability.Cost)
		if err != nil {
			panic("invalid mana cost in activated ability: " + err.Error())
		}
		zone := mtg.ZoneBattlefield
		if ability.Zone != "" {
			zone = ability.Zone
		}
		activatedAbilities = append(activatedAbilities, Ability{
			id:          fmt.Sprintf("%s-%d", config.ID, i+1),
			name:        ability.Name,
			cost:        abilityCost,
			effectSpecs: ability.EffectSpecs,
			zone:        zone,
			source:      permanent,
		})
	}
	permanent.activatedAbilities = activatedAbilities
	return permanent
}

func LoadCardFromConfig(playerID string, config gobtest.CardConfig) Card {
	card := Card{
		id:           config.ID,
		name:         config.Name,
		controller:   playerID,
		owner:        playerID,
		loyalty:      config.Loyalty,
		spellAbility: config.SpellAbility,
		rulesText:    config.RulesText,
		power:        config.Power,
		subtypes:     config.Subtypes,
	}
	manaCost, err := cost.ParseManaCost(config.ManaCost)
	if err != nil {
		panic("invalid mana cost in card config: " + err.Error())
	}
	card.manaCost = manaCost
	var ActivatedAbilities []Ability
	for i, ability := range config.ActivatedAbilities {
		abilityCost, err := cost.ParseManaCost(ability.Cost)
		if err != nil {
			panic("invalid mana cost in activated ability: " + err.Error())
		}
		zone := mtg.ZoneBattlefield
		if ability.Zone != "" {
			zone = ability.Zone
		}
		ActivatedAbilities = append(ActivatedAbilities, Ability{
			id:          fmt.Sprintf("%s-%d", config.ID, i+1),
			name:        ability.Name,
			cost:        abilityCost,
			effectSpecs: ability.EffectSpecs,
			zone:        zone,
			source:      card,
		})
	}
	card.activatedAbilities = ActivatedAbilities
	var staticAbilities []StaticAbility
	for _, ability := range config.StaticAbilities {
		abilityCost, err := cost.ParseManaCost(ability.Cost)
		if err != nil {
			panic("invalid mana cost in static ability: " + err.Error())
		}
		staticAbilities = append(staticAbilities, StaticAbility{
			name:      ability.Name,
			cost:      abilityCost,
			modifiers: ability.Modifiers,
		})
	}
	card.staticAbilities = staticAbilities
	return card
}

func LoadSpellFromConfig(config gobtest.SpellConfig) Spell {
	spell := Spell{
		id:         config.ID,
		name:       config.Name,
		controller: config.Controller,
		owner:      config.Owner,
		manaCost:   config.ManaCost,
	}
	return spell
}
