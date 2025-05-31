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

type Revealed struct {
	cards []gob.Card
}

func NewRevealed() Revealed {
	revealed := Revealed{
		cards: []gob.Card{},
	}
	return revealed
}

func (r Revealed) Add(card gob.Card) Revealed {
	return Revealed{cards: add.Item(r.cards, card)}
}

func (r Revealed) AddTop(c gob.Card) Revealed {
	return Revealed{
		cards: add.Item([]gob.Card{c}, r.cards...),
	}
}

func (r Revealed) Get(id string) (gob.Card, bool) {
	return query.Get(r.cards, id)
}

func (r Revealed) GetAll() []gob.Card {
	return query.GetAll(r.cards)
}

func (r Revealed) Name() string {
	return string(mtg.ZoneRevealed)
}

func (r Revealed) Remove(id string) (Revealed, bool) {
	cards, ok := remove.By(r.cards, has.ID(id))
	if !ok {
		return r, false
	}
	return Revealed{cards: cards}, true
}

func (r Revealed) Take(id string) (gob.Card, Revealed, bool) {
	card, cards, ok := take.By(r.cards, has.ID(id))
	if !ok {
		return gob.Card{}, r, false
	}
	return card, Revealed{cards: cards}, true
}

func (r Revealed) Size() int {
	return len(r.cards)
}

func (r Revealed) ZoneType() mtg.Zone {
	return mtg.ZoneRevealed
}
