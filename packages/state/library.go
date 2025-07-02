package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/add"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/remove"
	"deckronomicon/packages/query/take"
)

type Library struct {
	cards []*gob.Card
}

func NewLibrary(cards []*gob.Card) *Library {
	library := Library{
		cards: cards,
	}
	return &library
}

func (l *Library) Add(card *gob.Card) *Library {
	return &Library{
		cards: add.Item(l.cards, card),
	}
}

func (l *Library) AddTop(c *gob.Card) *Library {
	return &Library{
		cards: add.Item([]*gob.Card{c}, l.cards...),
	}
}

func (l *Library) Contains(predicate query.Predicate) bool {
	return query.Contains(l.cards, predicate)
}

func (l *Library) Find(predicate query.Predicate) (*gob.Card, bool) {
	return query.Find(l.cards, predicate)
}

func (l *Library) FindAll(predicate query.Predicate) []*gob.Card {
	return query.FindAll(l.cards, predicate)
}

func (l *Library) Get(id string) (*gob.Card, bool) {
	return query.Get(l.cards, id)
}

func (l *Library) GetAll() []*gob.Card {
	return query.GetAll(l.cards)
}

func (l *Library) GetN(n int) []*gob.Card {
	return query.GetN(l.cards, n)
}

func (l *Library) Name() string {
	return string(mtg.ZoneLibrary)
}

func (l *Library) Peek() *gob.Card {
	if len(l.cards) == 0 {
		return nil
	}
	return l.cards[0]
}

func (l *Library) Remove(id string) (*Library, bool) {
	cards, ok := remove.By(l.cards, has.ID(id))
	if !ok {
		return nil, false
	}
	return &Library{cards: cards}, true
}

func (l *Library) Take(id string) (*gob.Card, *Library, bool) {
	card, cards, ok := take.By(l.cards, has.ID(id))
	if !ok {
		return nil, nil, false
	}
	return card, &Library{cards: cards}, true
}

func (l *Library) TakeBy(query query.Predicate) (*gob.Card, *Library, bool) {
	taken, remaining, ok := take.By(l.cards, query)
	if !ok {
		return nil, nil, false
	}
	return taken, &Library{cards: remaining}, true
}

func (l *Library) TakeN(n int) ([]*gob.Card, *Library) {
	cards, remaining := take.N(l.cards, n)
	return cards, &Library{cards: remaining}
}

func (l *Library) TakeTop() (*gob.Card, *Library, bool) {
	card, cards, ok := take.Top(l.cards)
	if !ok {
		return nil, nil, false
	}
	return card, &Library{cards: cards}, true
}

func (l *Library) Size() int {
	return len(l.cards)
}

func (l *Library) ZoneType() mtg.Zone {
	return mtg.ZoneLibrary
}
