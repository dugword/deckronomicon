package state

// TODO: I think a lot of these methods should just be methods off of the View
// interface.
import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/take"
	"fmt"
	"math/rand/v2"
)

// Library represents the player's library.
type Library struct {
	cards []gob.Card
}

// NewLibrary creates a new Library instance.
func NewLibrary(deckList []gob.Card) Library {
	library := Library{
		cards: deckList,
	}
	return library
}

// Add adds a card to the bottom of the library.
func (l Library) Add(card gob.Card) {
	l.cards = append(l.cards, card)
}

// AddTop adds a card to the top of the library.
func (l Library) AddTop(c gob.Card) {
	l.cards = append([]gob.Card{c}, l.cards...)
}

func (l Library) Get(id string) (gob.Card, error) {
	for _, card := range l.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return gob.Card{}, fmt.Errorf("card with ID %s not found", id)
}

func (l Library) GetAll() []gob.Card {
	return l.cards
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

func (l Library) Remove(id string) error {
	for i, card := range l.cards {
		if card.ID() == id {
			l.cards = append(l.cards[:i], l.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with ID %s not found", id)
}

func (l Library) Take(id string) (gob.Card, error) {
	for i, card := range l.cards {
		if card.ID() != id {
			continue
		}
		l.cards = append(l.cards[:i], l.cards[i+1:]...)
		return card, nil
	}
	return gob.Card{}, fmt.Errorf("card with ID %s not found", id)
}

func (l Library) TakeBy(query query.Predicate) (gob.Card, error) {
	taken, remaining, err := take.By(l.cards, query)
	if err != nil {
		return gob.Card{}, fmt.Errorf("failed to take card by query: %w", err)
	}
	l.cards = remaining
	return taken, nil
}

func (l Library) Shift() (gob.Card, Library, bool) {
	if len(l.cards) == 0 {
		return gob.Card{}, l, false
	}
	taken := l.cards[0]
	newLibrary := Library{
		cards: l.cards[1:],
	}
	return taken, newLibrary, true
}

func (l Library) Size() int {
	return len(l.cards)
}

// Shuffle randomly shuffles the cards in the library.
func (l Library) Shuffle() {
	rand.Shuffle(len(l.cards), func(i, j int) {
		l.cards[i], l.cards[j] = l.cards[j], l.cards[i]
	})
}

func (l Library) ZoneType() mtg.Zone {
	return mtg.ZoneLibrary
}
