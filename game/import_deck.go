package game

import (
	"encoding/json"
	"io"
	"os"
)

type CardImport struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DeckImport struct {
	Cards []CardImport `json:"cards"`
}

func importDeck(filename string) (*Deck, error) {
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
	var deck Deck
	for _, cardImport := range deckImport.Cards {
		for range cardImport.Count {
			card := CardPool[cardImport.Name]
			deck.cards = append(deck.cards, &card)
		}
	}
	return &deck, err
}
