package game

import (
	"encoding/json"
	"fmt"
	"io"
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
func importDeck(filename string) (*Library, error) {
	var deckImport DeckImport
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &deckImport)
	if err != nil {
		return nil, err
	}
	library := NewLibrary()
	for _, cardImport := range deckImport.Cards {
		for range cardImport.Count {
			// TODO: maybe call NewCard here so each card get's a unique ID
			card, ok := CardPool[cardImport.Name]
			if !ok {
				return nil, fmt.Errorf("card %s not found in card pool", cardImport.Name)
			}
			library.PutOnTop(card)
		}
	}
	return library, err
}
