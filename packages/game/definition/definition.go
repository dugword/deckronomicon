package definition

// TODO Rename this package to not have an underscore, and maybe make it
// something else like raw card data or card definitions.

import (
	"deckronomicon/packages/game/mtg"
)

type EffectSpec struct {
	Name      string         `json:"Name,omitempty" yaml:"Name,omitempty"`
	Modifiers map[string]any `json:"Modifiers,omitempty" yaml:"Modifiers,omitempty"`
}

type ActivatedAbilitySpec struct {
	Name        string       `json:"Name,omitempty" yaml:"Name,omitempty"`
	Cost        string       `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	EffectSpecs []EffectSpec `json:"Effects,omitempty" yaml:"Effects,omitempty"`
	Speed       mtg.Speed    `json:"Speed,omitempty" yaml:"Speed,omitempty"`
	Zone        mtg.Zone     `json:"Zone,omitempty" yaml:"Zone,omitempty"`
}

// StaticAbility represents the specification of static ability.
type StaticAbilitySpec struct {
	Name      mtg.StaticKeyword `json:"Name,omitempty" yaml:"Name,omitempty"`
	Cost      string            `json:"Cost,omitempty" yaml:"Cost,omitempty"`
	Modifiers map[string]any    `json:"Modifiers,omitempty" yaml:"Modifiers,omitempty"`
}

type TriggeredAbilitySpec struct {
	EffectSpec []EffectSpec `json:"Effects,omitempty" yaml:"Effects,omitempty"`
}

// Card represents the underlying data structure for a card in the game.
type Card struct {
	ActivatedAbilitySpecs []ActivatedAbilitySpec `json:"ActivatedAbilities,omitempty" yaml:"ActivatedAbilities,omitempty"`
	CardTypes             []mtg.CardType         `json:"CardTypes,omitempty" yaml:"CardTypes,omitempty"`
	Colors                []string               `json:"Color,omitempty" yaml:"Color,omitempty"`
	Loyalty               int                    `json:"Loyalty,omitempty" yaml:"Loyalty,omitempty"`
	ManaCost              string                 `json:"ManaCost,omitempty" yaml:"ManaCost,omitempty"`
	Name                  string                 `json:"Name,omitempty" yaml:"Name,omitempty"`
	Power                 int                    `json:"Power,omitempty" yaml:"Power,omitempty"`
	RulesText             string                 `json:"RulesText,omitempty" yaml:"RulesText,omitempty"`
	SpellAbilitySpec      []EffectSpec           `json:"Effects,omitempty" yaml:"Effects,omitempty"`
	StaticAbilitySpecs    []StaticAbilitySpec    `json:"StaticAbilities,omitempty" yaml:"StaticAbilities,omitempty"`
	TriggeredAbilitySpecs []TriggeredAbilitySpec `json:"TriggeredAbilities,omitempty" yaml:"TriggeredAbilities,omitempty"`
	Subtypes              []mtg.Subtype          `json:"Subtypes,omitempty" yaml:"Subtypes,omitempty"`
	Supertypes            []mtg.Supertype        `json:"Supertypes,omitempty" yaml:"Supertypes,omitempty"`
	Toughness             int                    `json:"Toughness,omitempty" yaml:"Toughness,omitempty"`
}
