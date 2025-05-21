package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadCardsFromCardPool loads cards from a directory containing JSON files.
// TODO: I think this doesn't need to exist, and we should load cards during
// the deck loading phase, this would allow us to call NewCardFromImport for
// each card and generate a unique ID for each card.
func LoadCardPool(path string) (map[string]*Card, error) {
	cards := map[string]*Card{}
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		var cardImport CardImport
		// TODO, or can I just add the file here...
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&cardImport); err != nil {
			return fmt.Errorf("failed to unmarshal file %s: %w", path, err)
		}
		card, err := NewCardFromImport(cardImport)
		if err != nil {
			return fmt.Errorf("failed to create card from import: %w", err)
		}
		if _, ok := cards[card.Name()]; ok {
			return fmt.Errorf("duplicate card name detected: %s", card.Name)
		}
		cards[card.Name()] = card
		return nil
	}
	if err := filepath.Walk(path, walkFunc); err != nil {
		return nil, fmt.Errorf("could not load cards: %w", err)
	}
	return cards, nil
}

// MustLoadCardPool loads cards from a directory containing JSON files and
// panics on error.
// TODO: Only here because I don't want to deal with errors yet, this should
// be managed during the GameState initialization when the deck is being
// loaded.
func MustLoadCardPool(path string) map[string]*Card {
	cards, err := LoadCardPool(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load card pool: %v", err))
	}
	return cards
}

// CardPool is a map of card names to Card objects.
// TODO: This is a temporary hack and should be moved to the GameState
// initialization function.
var CardPool = MustLoadCardPool("./cards")
