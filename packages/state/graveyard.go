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

func (g Graveyard) Add(cards ...gob.Card) Graveyard {
	return Graveyard{
		cards: add.Item(g.cards, cards...),
	}
}

func (g Graveyard) AddTop(c gob.Card) Graveyard {
	return Graveyard{
		cards: add.Item([]gob.Card{c}, g.cards...),
	}
}

func (g Graveyard) Contains(predicate query.Predicate) bool {
	return query.Contains(g.cards, predicate)
}

func (g Graveyard) Get(id string) (gob.Card, bool) {
	return query.Get(g.cards, id)
}

func (g Graveyard) GetAll() []gob.Card {
	return query.GetAll(g.cards)
}

func (g Graveyard) Name() string {
	return string(mtg.ZoneGraveyard)
}

func (g Graveyard) Remove(id string) (Graveyard, bool) {
	cards, ok := remove.By(g.cards, has.ID(id))
	if !ok {
		return g, false
	}
	return Graveyard{cards: cards}, true
}

func (g Graveyard) Take(id string) (gob.Card, Graveyard, bool) {
	card, cards, ok := take.By(g.cards, has.ID(id))
	if !ok {
		return gob.Card{}, g, false
	}
	return card, Graveyard{cards: cards}, true
}

func (g Graveyard) TakeBy(predicate query.Predicate) (gob.Card, Graveyard, bool) {
	card, cards, ok := take.By(g.cards, predicate)
	if !ok {
		return gob.Card{}, g, false
	}
	return card, Graveyard{cards: cards}, true
}

func (g Graveyard) Size() int {
	return len(g.cards)
}

func (g Graveyard) ZoneType() mtg.Zone {
	return mtg.ZoneGraveyard
}
