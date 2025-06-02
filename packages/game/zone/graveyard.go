package zone

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/object"
	"fmt"
)

type Graveyard struct {
	cards []*object.Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() *Graveyard {
	return &Graveyard{
		cards: []*object.Card{},
	}
}

func (g *Graveyard) Add(card *object.Card) {
	g.cards = append(g.cards, card)
}

func (g *Graveyard) Get(id string) (*object.Card, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) GetAll() []*object.Card {
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

func (g *Graveyard) Take(id string) (*object.Card, error) {
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
