package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Revealed struct {
	cards []gob.Card
}

func NewRevealed() Revealed {
	revealed := Revealed{
		cards: []gob.Card{},
	}
	return revealed
}

func (r Revealed) Add(card gob.Card) {
	r.cards = append(r.cards, card)
}

func (r Revealed) Clear() {
	r.cards = []gob.Card{}
}

func (r Revealed) Get(id string) (gob.Card, error) {
	for _, card := range r.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return gob.Card{}, fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r Revealed) GetAll() []gob.Card {
	return r.cards
}

func (r Revealed) Name() string {
	return string(mtg.ZoneRevealed)
}

func (r Revealed) Remove(id string) error {
	for i, card := range r.cards {
		if card.ID() == id {
			r.cards = append(r.cards[:i], r.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r Revealed) Take(id string) (gob.Card, error) {
	for i, card := range r.cards {
		if card.ID() == id {
			r.cards = append(r.cards[:i], r.cards[i+1:]...)
			return card, nil
		}
	}
	return gob.Card{}, fmt.Errorf("card with id %s not found in revealed zone", id)
}

func (r Revealed) Size() int {
	return len(r.cards)
}

func (r Revealed) ZoneType() mtg.Zone {
	return mtg.ZoneRevealed
}
