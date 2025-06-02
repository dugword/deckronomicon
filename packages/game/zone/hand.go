package zone

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/object"
	"fmt"
)

// Hand represents a player's hand of cards.
type Hand struct {
	cards []*object.Card
}

// NewHand creates a new Hand instance.
func NewHand() *Hand {
	return &Hand{
		cards: []*object.Card{},
	}
}

func (h *Hand) Add(card *object.Card) {
	h.cards = append(h.cards, card)
}

func (h *Hand) Get(id string) (*object.Card, error) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) GetAll() []*object.Card {
	return h.cards
}

// TODO: think if I want this to be "%s's Hand" or just "Hand"
// Right now this is for the choose.Source interface.
func (h *Hand) Name() string {
	return string(mtg.ZoneHand)
}

func (h *Hand) Remove(id string) error {
	for i, card := range h.cards {
		if card.ID() == id {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) Take(id string) (*object.Card, error) {
	for i, card := range h.cards {
		if card.ID() == id {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) Size() int {
	return len(h.cards)
}

func (h *Hand) ZoneType() mtg.Zone {
	return mtg.ZoneHand
}
