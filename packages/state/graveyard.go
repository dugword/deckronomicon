package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Graveyard struct {
	cards []gob.Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() Graveyard {
	graveyard := Graveyard{
		cards: []gob.Card{},
	}
	return graveyard
}

func (g Graveyard) Append(cards ...gob.Card) Graveyard {
	newCards := append(g.cards[:], cards...)
	return Graveyard{
		cards: newCards,
	}
}

func (g Graveyard) Get(id string) (gob.Card, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return gob.Card{}, fmt.Errorf("card witg ID %s not found", id)
}

func (g Graveyard) GetAll() []gob.Card {
	return g.cards
}

func (g Graveyard) Name() string {
	return string(mtg.ZoneGraveyard)
}

func (g Graveyard) Remove(id string) error {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card witg ID %s not found", id)
}

func (g Graveyard) Take(id string) (gob.Card, error) {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return card, nil
		}
	}
	return gob.Card{}, fmt.Errorf("card witg ID %s not found", id)
}

func (g Graveyard) Size() int {
	return len(g.cards)
}

func (g Graveyard) ZoneType() mtg.Zone {
	return mtg.ZoneGraveyard
}
