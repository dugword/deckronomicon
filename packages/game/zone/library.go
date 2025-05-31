package zone

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/take"
	"fmt"
	"math/rand/v2"
)

// TODO should this live here?
type state interface {
	GetNextID() string
}

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

func BuildLibrary(
	state state,
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

/*
// TODO not sure if this is safe to modify a slice as I iterate over it
	var taken []*card.Card
	for i, card := range l.cards {
		if query(card) {
			l.cards = append(l.cards[:i], l.cards[i+1:]...)
		}
	}
	return taken

*/

// TakeCards removes the top N cards from the library and returns them.
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

// Library Specific Methods

// TODO move these to player?

/*
// AvailableActivatedAbilities returns a list of activated abilities that can
// be activated from the library. This exits to satisfy the Zone interface.
// Cards in the library cannot have a activated ability.
func (l *Library) AvailableActivatedAbilities(*GameState, *Player) []game.Object {
	return nil
}

// AvailableToPlay returns a list of cards that can be played from the
// library. This exists to satisfy the Zone interface. Cards in the library
// generally cannot be played. (Some cards enable this, but they are
// exceptions and not currently implemented.)
func (l *Library) AvailableToPlay(*GameState, *Player) []game.Object {
	return nil
}
*/
