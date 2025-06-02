package zone

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Graveyard struct {
	cards []*card.Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() *Graveyard {
	return &Graveyard{
		cards: []*card.Card{},
	}
}

func (g *Graveyard) Add(card *card.Card) {
	g.cards = append(g.cards, card)
}

func (g *Graveyard) Get(id string) (*card.Card, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) GetAll() []*card.Card {
	return g.cards
}

func (g *Graveyard) Name() string {
	return string(mtg.ZoneGraveyard)
}

func (g *Graveyard) Remove(id string) error {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) Take(id string) (*card.Card, error) {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) Size() int {
	return len(g.cards)
}

func (g *Graveyard) ZoneType() mtg.Zone {
	return mtg.ZoneGraveyard
}
