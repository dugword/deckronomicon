package definition

// TODO Rename this package to not have an underscore, and maybe make it
// something else like raw card data or card definitions.

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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

// LoadCardDefinitions walks the provided path and loads all YAML files into a map of
// card names to card data. It returns an error if any file cannot be read
// or if there are duplicate card names.
func LoadCardDefinitions(path string) (map[string]Card, error) {
	definitions := map[string]Card{}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read card file at %q: %w", path, err)
		}
		var card Card
		if err := yaml.Unmarshal(data, &card); err != nil {
			return fmt.Errorf("failed to unmarshal card data in %q: %w", path, err)
		}
		if card.Name == "" {
			return fmt.Errorf("card in %q is missing a name", path)
		}
		if _, ok := definitions[card.Name]; ok {
			return fmt.Errorf("duplicate card name detected in %q: %s", path, card.Name)
		}
		definitions[card.Name] = card
		return nil
	}
	if err := filepath.Walk(path, walkFunc); err != nil {
		return nil, fmt.Errorf("failed to load card definitions in %q: %w", path, err)
	}
	return definitions, nil
}
