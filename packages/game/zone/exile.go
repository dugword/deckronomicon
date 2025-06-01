package zone

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"errors"
)

type Exile struct {
	cards []*card.Card
}

func NewExile() *Exile {
	return &Exile{
		cards: []*card.Card{},
	}
}

func (e *Exile) Add(card *card.Card) {
	e.cards = append(e.cards, card)
}

func (e *Exile) Get(id string) (*card.Card, error) {
	for _, card := range e.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}

func (e *Exile) GetAll() []*card.Card {
	var cards []*card.Card
	for _, card := range e.cards {
		cards = append(cards, card)
	}
	return cards
}

func (e *Exile) Name() string {
	return string(mtg.ZoneExile)
}

func (e *Exile) Remove(id string) error {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return nil
		}
	}
	return errors.New("card not found in exile")
}
func (e *Exile) Take(id string) (*card.Card, error) {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}
func (e *Exile) Size() int {
	return len(e.cards)
}
func (e *Exile) ZoneType() mtg.Zone {
	return mtg.ZoneExile
}
