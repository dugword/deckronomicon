package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EffectSpec struct {
	ID        string           `json:"ID,omitempty"`
	Modifiers []EffectModifier `json:"Modifiers,omitempty"`
}

type EffectModifier struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

type ActivatedAbilitySpec struct {
	Cost        string       `json:"Cost,omitempty"`
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	// TODO: This might need to be a slice if an ability is activatedable from
	// multiple zones.
	Zone string `json:"Zone,omitempty"`
}

type SpellAbilitySpec struct {
	// Cost // TODO: AdditionalCosts?
	EffectSpecs []EffectSpec `json:"Effects,omitempty"`
	Zone        string       `json:"Zone,omitempty"`
}

// StaticAbility represents the specification of static ability.
type StaticAbilitySpec struct {
	ID        string           `json:"ID,omitempty"`
	Modifiers []EffectModifier `json:"Modifiers,omitempty"`
}

type TriggeredAbilitySpec struct {
	// Cost Cost // TODO: Additoonal Cost
	EffectSpec []EffectSpec `json:"Effects,omitempty"`
}

// CardData represents the underlying data structure for a card in the game.
type CardData struct {
	ActivatedAbilitySpecs []*ActivatedAbilitySpec `json:"ActivatedAbilities,omitempty"`
	CardTypes             []string                `json:"CardTypes,omitempty"`
	Colors                []string                `json:"Color,omitempty"`
	Loyalty               int                     `json:"Loyalty,omitempty"`
	ManaCost              string                  `json:"ManaCost,omitempty"`
	Name                  string                  `json:"Name,omitempty"`
	Power                 int                     `json:"Power,omitempty"`
	RulesText             string                  `json:"RulesText,omitempty"`
	SpellAbilitySpec      *SpellAbilitySpec       `json:"SpellAbility,omitempty"`
	StaticAbilitySpecs    []*StaticAbilitySpec    `json:"StaticAbilities,omitempty"`
	TriggeredAbilitySpecs []*TriggeredAbilitySpec `json:"TriggeredAbilities,omitempty"`
	Subtypes              []string                `json:"Subtypes,omitempty"`
	Supertypes            []string                `json:"Supertypes,omitempty"`
	Toughness             int                     `json:"Toughness,omitempty"`
}

// LoadCardPoolData walks the provided path and loads all JSON files into a map of
// card names to card data. It returns an error if any file cannot be read
// or if there are duplicate card names.
func LoadCardPoolData(path string) (map[string]CardData, error) {
	cardPoolData := map[string]CardData{}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read card file %s: %w", path, err)
		}
		var cardData CardData
		if err := json.Unmarshal(data, &cardData); err != nil {
			return fmt.Errorf("failed to unmarshal card data %s: %w", path, err)
		}
		if _, ok := cardPoolData[cardData.Name]; ok {
			return fmt.Errorf("duplicate card name detected: %s", cardData.Name)
		}
		cardPoolData[cardData.Name] = cardData
		return nil
	}
	if err := filepath.Walk(path, walkFunc); err != nil {
		return nil, fmt.Errorf("could not load cards: %w", err)
	}
	return cardPoolData, nil
}
