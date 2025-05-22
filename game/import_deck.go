package game

import (
	"encoding/json"
	"fmt"
	"os"
)

// DeckImportCard represents a card in the deck import file.
type DeckImportCard struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// DeckImport represents the structure of the deck import file.
type DeckImport struct {
	Cards []DeckImportCard `json:"cards"`
}

// importDeck imports a deck from a JSON file. The file should contain a list of
// cards and their counts. The function returns a Library containing the
// imported cards or an error if the import fails.
func importDeck(filename string, cardPool string) (*Library, error) {
	var deckImport DeckImport
	deckData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read deck import file: %w", err)
	}
	if err := json.Unmarshal(deckData, &deckImport); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck import: %w", err)
	}
	cardPoolData, err := LoadCardPoolData(cardPool)
	if err != nil {
		return nil, fmt.Errorf("failed to load card pool data: %w", err)
	}
	library := NewLibrary()
	for _, deckImportCard := range deckImport.Cards {
		for range deckImportCard.Count {
			cardData, ok := cardPoolData[deckImportCard.Name]
			if !ok {
				return nil, fmt.Errorf("card %s not found in card pool data", deckImportCard.Name)
			}
			card, err := NewCardFromCardData(cardData)
			if err != nil {
				return nil, fmt.Errorf("failed to create card %s: %w", deckImportCard.Name, err)
			}
			library.PutOnTop(card)
		}
	}
	return library, err
}
