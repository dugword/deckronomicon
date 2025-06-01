package zone

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/take"
	"fmt"
	"math/rand/v2"
)

// Library represents the player's library.
type Library struct {
	cards []*card.Card
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	library := Library{
		cards: []*card.Card{},
	}
	return &library
}

// TODO Not sure if I like this here.
func BuildLibrary(
	state core.State,
	deckList configs.DeckList,
	cardDefinitions map[string]definition.Card,
) (*Library, error) {
	library := NewLibrary()
	for _, entry := range deckList.Cards {
		for range entry.Count {
			cardDefinition, ok := cardDefinitions[entry.Name]
			if !ok {
				return nil, fmt.Errorf(
					"card %s not found in card definitions",
					entry.Name,
				)
			}
			c, err := card.NewCardFromCardDefinition(state, cardDefinition)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to create c %s: %w",
					entry.Name,
					err,
				)
			}
			library.AddTop(c)
		}
	}
	return library, nil
}

// Add adds a card to the bottom of the library.
func (l *Library) Add(card *card.Card) {
	l.cards = append(l.cards, card)
}

// AddTop adds a card to the top of the library.
func (l *Library) AddTop(c *card.Card) {
	l.cards = append([]*card.Card{c}, l.cards...)
}

func (l *Library) Get(id string) (*card.Card, error) {
	for _, card := range l.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (l *Library) GetAll() []*card.Card {
	return l.cards
}

func (l *Library) Name() string {
	return string(mtg.ZoneLibrary)
}

// Peek returns the top N cards without modifying the library.
func (l *Library) Peek() *card.Card {
	if len(l.cards) == 0 {
		return nil
	}
	return l.cards[0]
}

func (l *Library) Remove(id string) error {
	for i, card := range l.cards {
		if card.ID() == id {
			l.cards = append(l.cards[:i], l.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with ID %s not found", id)
}

func (l *Library) Take(id string) (*card.Card, error) {
	for i, card := range l.cards {
		if card.ID() != id {
			continue
		}
		l.cards = append(l.cards[:i], l.cards[i+1:]...)
		return card, nil
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (l *Library) TakeBy(query query.Predicate) (*card.Card, error) {
	taken, remaining, err := take.By(l.cards, query)
	if err != nil {
		return nil, fmt.Errorf("failed to take card by query: %w", err)
	}
	l.cards = remaining
	return taken, nil
}

func (l *Library) TakeTop() (*card.Card, error) {
	if len(l.cards) == 0 {
		return nil, mtg.ErrLibraryEmpty
	}
	taken := l.cards[0]
	l.cards = l.cards[1:]
	return taken, nil
}

func (l *Library) Size() int {
	return len(l.cards)
}

// Shuffle randomly shuffles the cards in the library.
func (l *Library) Shuffle() {
	rand.Shuffle(len(l.cards), func(i, j int) {
		l.cards[i], l.cards[j] = l.cards[j], l.cards[i]
	})
}

func (l *Library) ZoneType() mtg.Zone {
	return mtg.ZoneLibrary
}
