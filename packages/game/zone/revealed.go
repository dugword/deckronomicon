package zone

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Revealed struct {
	cards []*card.Card
}

func NewRevealed() *Revealed {
	return &Revealed{
		cards: []*card.Card{},
	}
}

func (r *Revealed) Add(card *card.Card) {
	r.cards = append(r.cards, card)
}

func (r *Revealed) Clear() {
	r.cards = []*card.Card{}
}

func (r *Revealed) Get(id string) (*card.Card, error) {
	for _, card := range r.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r *Revealed) GetAll() []*card.Card {
	return r.cards
}

func (r *Revealed) Name() string {
	return string(mtg.ZoneRevealed)
}

func (r *Revealed) Remove(id string) error {
	for i, card := range r.cards {
		if card.ID() == id {
			r.cards = append(r.cards[:i], r.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r *Revealed) Take(id string) (*card.Card, error) {
	for i, card := range r.cards {
		if card.ID() == id {
			r.cards = append(r.cards[:i], r.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r *Revealed) Size() int {
	return len(r.cards)
}

func (r *Revealed) ZoneType() mtg.Zone {
	return mtg.ZoneRevealed
}
