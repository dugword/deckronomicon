package state

// TODO: I think a lot of these methods should just be methods off of the View
// interface.
import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/add"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/remove"
	"deckronomicon/packages/query/take"
)

// Library represents the player's library.
type Library struct {
	cards []gob.Card
}

// NewLibrary creates a new Library instance.
func NewLibrary(cards []gob.Card) Library {
	library := Library{
		cards: cards,
	}
	return library
}

// Add adds a card to the bottom of the library.
func (l Library) Add(card gob.Card) Library {
	return Library{
		cards: add.Item(l.cards, card),
	}
}

// AddTop adds a card to the top of the library.
func (l Library) AddTop(c gob.Card) Library {
	return Library{
		cards: add.Item([]gob.Card{c}, l.cards...),
	}
}

func (l Library) Get(id string) (gob.Card, bool) {
	return query.Get(l.cards, id)
}

func (l Library) GetAll() []gob.Card {
	return query.GetAll(l.cards)
}

func (l Library) Name() string {
	return string(mtg.ZoneLibrary)
}

// Peek returns the top N cards without modifying the library.
func (l Library) Peek() gob.Card {
	if len(l.cards) == 0 {
		return gob.Card{}
	}
	return l.cards[0]
}

func (l Library) Remove(id string) (Library, bool) {
	cards, ok := remove.By(l.cards, has.ID(id))
	if !ok {
		return l, false
	}
	return Library{cards: cards}, true
}

func (l Library) Take(id string) (gob.Card, Library, bool) {
	card, cards, ok := take.By(l.cards, has.ID(id))
	if !ok {
		return gob.Card{}, l, false
	}
	return card, Library{cards: cards}, true
}

func (l Library) TakeBy(query query.Predicate) (gob.Card, Library, bool) {
	taken, remaining, ok := take.By(l.cards, query)
	if !ok {
		return gob.Card{}, l, false
	}
	return taken, Library{cards: remaining}, true
}

func (l Library) TakeTop() (gob.Card, Library, bool) {
	card, cards, ok := take.Top(l.cards)
	if !ok {
		return gob.Card{}, l, false
	}
	return card, Library{cards: cards}, true
}

func (l Library) Size() int {
	return len(l.cards)
}

func (l Library) ZoneType() mtg.Zone {
	return mtg.ZoneLibrary
}
