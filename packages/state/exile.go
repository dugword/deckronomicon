package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"errors"
)

type Exile struct {
	cards []gob.Card
}

func NewExile() Exile {
	exile := Exile{
		cards: []gob.Card{},
	}
	return exile
}

func (e Exile) Add(card gob.Card) {
	e.cards = append(e.cards, card)
}

func (e Exile) Get(id string) (gob.Card, error) {
	for _, card := range e.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return gob.Card{}, errors.New("card not found in exile")
}

func (e Exile) GetAll() []gob.Card {
	var cards []gob.Card
	for _, card := range e.cards {
		cards = append(cards, card)
	}
	return cards
}

func (e Exile) Name() string {
	return string(mtg.ZoneExile)
}

func (e Exile) Remove(id string) error {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return nil
		}
	}
	return errors.New("card not found in exile")
}
func (e Exile) Take(id string) (gob.Card, error) {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return card, nil
		}
	}
	return gob.Card{}, errors.New("card not found in exile")
}
func (e Exile) Size() int {
	return len(e.cards)
}
func (e Exile) ZoneType() mtg.Zone {
	return mtg.ZoneExile
}
