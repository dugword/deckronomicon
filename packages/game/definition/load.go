package definition

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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
		card, err := LoadCardDefinition(path)
		if err != nil {
			return err
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

func LoadCardDefinition(path string) (Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Card{}, fmt.Errorf("failed to read card file at %q: %w", path, err)
	}
	var card Card
	if err := yaml.Unmarshal(data, &card); err != nil {
		return Card{}, fmt.Errorf("failed to unmarshal card data in %q: %w", path, err)
	}
	if card.Name == "" {
		return Card{}, fmt.Errorf("card in %q is missing a name", path)
	}
	return card, nil
}
