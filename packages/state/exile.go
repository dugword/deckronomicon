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

type Exile struct {
	cards []*gob.Card
}

func NewExile() *Exile {
	exile := Exile{
		cards: []*gob.Card{},
	}
	return &exile
}

func (e *Exile) Add(cards ...*gob.Card) *Exile {
	return &Exile{
		cards: add.Item(e.cards, cards...),
	}
}

func (e *Exile) AddTop(c *gob.Card) *Exile {
	return &Exile{
		cards: add.Item([]*gob.Card{c}, e.cards...),
	}
}

func (e *Exile) Contains(predicate query.Predicate) bool {
	return query.Contains(e.cards, predicate)
}

func (e *Exile) Get(id string) (*gob.Card, bool) {
	return query.Get(e.cards, id)
}

func (e *Exile) GetAll() []*gob.Card {
	return query.GetAll(e.cards)
}

func (e *Exile) Name() string {
	return string(mtg.ZoneExile)
}

func (e *Exile) Remove(id string) (*Exile, bool) {
	cards, ok := remove.By(e.cards, has.ID(id))
	if !ok {
		return nil, false
	}
	return &Exile{cards: cards}, true
}

func (e *Exile) Take(id string) (*gob.Card, *Exile, bool) {
	card, cards, ok := take.By(e.cards, has.ID(id))
	if !ok {
		return nil, nil, false
	}
	return card, &Exile{cards: cards}, true
}

func (e *Exile) TakeBy(predicate query.Predicate) (*gob.Card, *Exile, bool) {
	card, cards, ok := take.By(e.cards, predicate)
	if !ok {
		return nil, nil, false
	}
	return card, &Exile{cards: cards}, true
}

func (e *Exile) Size() int {
	return len(e.cards)
}

func (e *Exile) ZoneType() mtg.Zone {
	return mtg.ZoneExile
}
