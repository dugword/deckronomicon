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

// Hand represents a player's hand of cards.
type Hand struct {
	cards []gob.Card
}

// NewHand creates a new Hand instance.
func NewHand() Hand {
	return Hand{
		cards: []gob.Card{},
	}
}

func (h Hand) Add(cards ...gob.Card) Hand {
	return Hand{
		cards: add.Item(h.cards, cards...),
	}
}

func (h Hand) AddTop(c gob.Card) Hand {
	return Hand{
		cards: add.Item([]gob.Card{c}, h.cards...),
	}
}

func (h Hand) Contains(predicate query.Predicate) bool {
	return query.Contains(h.cards, predicate)
}

func (h Hand) Find(predicate query.Predicate) (gob.Card, bool) {
	return query.Find(h.cards, predicate)
}

func (h Hand) Get(id string) (gob.Card, bool) {
	return query.Get(h.cards, id)
}

func (h Hand) GetAll() []gob.Card {
	return query.GetAll(h.cards)
}

// TODO: think if I want this to be "%s's Hand" or just "Hand"
// Right now this is for the choose.Source interface.
func (h Hand) Name() string {
	return string(mtg.ZoneHand)
}

func (h Hand) Remove(id string) (Hand, bool) {
	cards, ok := remove.By(h.cards, has.ID(id))
	if !ok {
		return h, false
	}
	return Hand{cards: cards}, true
}

func (h Hand) Size() int {
	return len(h.cards)
}

func (h Hand) Take(id string) (gob.Card, Hand, bool) {
	card, cards, ok := take.By(h.cards, has.ID(id))
	if !ok {
		return gob.Card{}, h, false
	}
	return card, Hand{cards: cards}, true
}

func (h Hand) ZoneType() mtg.Zone {
	return mtg.ZoneHand
}
