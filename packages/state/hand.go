package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
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

func (h Hand) Append(cards ...gob.Card) Hand {
	newCards := append(h.cards[:], cards...)
	return Hand{
		cards: newCards,
	}
}

func (h Hand) Contains(predicate query.Predicate) bool {
	return query.Contains(h.cards, predicate)
}

func (h Hand) Find(predicate query.Predicate) (gob.Card, bool) {
	return query.Find(h.cards, predicate)
}

func (h Hand) Get(id string) (gob.Card, bool) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, true
		}
	}
	return gob.Card{}, false
}

func (h Hand) GetAll() []gob.Card {
	return h.cards
}

// TODO: think if I want this to be "%s's Hand" or just "Hand"
// Right now this is for the choose.Source interface.
func (h Hand) Name() string {
	return string(mtg.ZoneHand)
}

func (h Hand) Remove(id string) error {
	for i, card := range h.cards {
		if card.ID() == id {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card with ID %s not found", id)
}

func (h Hand) Take(id string) (gob.Card, Hand, bool) {
	remaining := Hand{}
	taken := gob.Card{}
	for _, card := range h.cards {
		if card.ID() == id {
			taken = card
			continue
		}
		remaining.cards = append(remaining.cards, card)
	}
	if taken.ID() == "" {
		return gob.Card{}, h, false
	}
	return taken, remaining, true
}

func (h Hand) Size() int {
	return len(h.cards)
}

func (h Hand) ZoneType() mtg.Zone {
	return mtg.ZoneHand
}
