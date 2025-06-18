package definition

// TODO Rename this package to not have an underscore, and maybe make it
// something else like raw card data or card definitions.

import (
	"deckronomicon/packages/game/mtg"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EffectSpec struct {
	Name      string          `json:"Name,omitempty"`
	Modifiers json.RawMessage `json:"Modifiers,omitempty"`
}

type ActivatedAbilitySpec struct {
	Name        string       `json:"Name,omitempty"`
	Cost        string       `json:"Cost,omitempty"`
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	// TODO: This might need to be a slice if an ability is activatedable from
	// multiple zones.
	Speed string `json:"Speed,omitempty"`
	Zone  string `json:"Zone,omitempty"`
}

type SpellAbilitySpec struct {
	// Cost // TODO: AdditionalCosts?
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	Zone        string       `json:"Zone,omitempty"`
}

// StaticAbility represents the specification of static ability.
type StaticAbilitySpec struct {
	Name      mtg.StaticKeyword `json:"Name,omitempty"`
	Cost      string            `json:"Cost,omitempty"`
	Modifiers json.RawMessage   `json:"Modifiers,omitempty"`
}

type TriggeredAbilitySpec struct {
	// Cost Cost // TODO: Additoonal Cost
	EffectSpec []EffectSpec `json:"Effects,omitempty"`
}

// Card represents the underlying data structure for a card in the game.
type Card struct {
	ActivatedAbilitySpecs []ActivatedAbilitySpec `json:"ActivatedAbilities,omitempty"`
	CardTypes             []mtg.CardType         `json:"CardTypes,omitempty"`
	Colors                []string               `json:"Color,omitempty"`
	Loyalty               int                    `json:"Loyalty,omitempty"`
	ManaCost              string                 `json:"ManaCost,omitempty"`
	Name                  string                 `json:"Name,omitempty"`
	Power                 int                    `json:"Power,omitempty"`
	RulesText             string                 `json:"RulesText,omitempty"`
	SpellAbilitySpec      SpellAbilitySpec       `json:"SpellAbility,omitempty"`
	StaticAbilitySpecs    []StaticAbilitySpec    `json:"StaticAbilities,omitempty"`
	TriggeredAbilitySpecs []TriggeredAbilitySpec `json:"TriggeredAbilities,omitempty"`
	Subtypes              []mtg.Subtype          `json:"Subtypes,omitempty"`
	Supertypes            []mtg.Supertype        `json:"Supertypes,omitempty"`
	Toughness             int                    `json:"Toughness,omitempty"`
}

// LoadCardDefinitions walks the provided path and loads all JSON files into a map of
// card names to card data. It returns an error if any file cannot be read
// or if there are duplicate card names.
func LoadCardDefinitions(path string) (map[string]Card, error) {
	definitions := map[string]Card{}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read card file at %q: %w", path, err)
		}
		var card Card
		if err := json.Unmarshal(data, &card); err != nil {
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
