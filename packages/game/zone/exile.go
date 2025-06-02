package zone

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/object"
	"errors"
)

type Exile struct {
	cards []*object.Card
}

func NewExile() *Exile {
	return &Exile{
		cards: []*object.Card{},
	}
}

func (e *Exile) Add(card *object.Card) {
	e.cards = append(e.cards, card)
}

func (e *Exile) Get(id string) (*object.Card, error) {
	for _, card := range e.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}

func (e *Exile) GetAll() []*object.Card {
	var cards []*object.Card
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
func (e *Exile) Take(id string) (*object.Card, error) {
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
