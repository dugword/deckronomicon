package game

import (
	"math/rand/v2"
)

// Library represents the player's library.
type Library struct {
	cards []*Card
}

// Cards returns the cards in the hand.
// TODO: for display remove later - why? Probably should directly manipulate
// this
func (l *Library) Cards() []*Card {
	return l.cards
}

// NewLibrary creates a new Library instance.
func NewLibrary() *Library {
	library := Library{
		cards: []*Card{},
	}
	return &library
}

// ChoseCardsBySubtype returns a list of cards of the specified subtype.
func (l *Library) ChooseCardsBySubtype(subtype Subtype) []Choice {
	var choices []Choice
	for i, card := range l.cards {
		if card.HasSubtype(subtype) {
			choices = append(choices, Choice{Name: card.Name(), Index: i})
		}
	}
	return choices
}

// Shuffle randomly shuffles the cards in the library.
func (l *Library) Shuffle() {
	rand.Shuffle(len(l.cards), func(i, j int) {
		l.cards[i], l.cards[j] = l.cards[j], l.cards[i]
	})
}

// TakeCardByIndex removes a card from the library at the specified index and
// returns it.
// TODO: Replace this with a way to select a card by ID, seems more robust.
// TODO: Standardize with hand/battlefield/graveyard take/get/find methods.
func (l *Library) TakeCardByIndex(index int) (*Card, error) {
	if index < 0 || index >= len(l.cards) {
		return nil, ErrLibraryEmpty
	}
	card := l.cards[index]
	l.cards = append(l.cards[:index], l.cards[index+1:]...)
	return card, nil
}

// TakeCards removes the top N cards from the library and returns them.
func (l *Library) TakeCards(n int) ([]*Card, error) {
	// TODO: I don't like this error handling
	var err error
	if n <= 0 || len(l.cards) == 0 {
		return l.cards, ErrLibraryEmpty
	}
	if n > len(l.cards) {
		n = len(l.cards)
		err = ErrLibraryEmpty
	}
	taken, remaining := l.cards[:n], l.cards[n:]
	l.cards = remaining
	return taken, err
}

// Peek returns the top N cards without modifying the library.
func (l *Library) Peek(pile []*Card, n int) []*Card {
	if n > len(pile) {
		n = len(pile)
	}
	return pile[:n]
}

// PutOnTop places the specified cards on top of the library.
func (l *Library) PutOnTop(cards ...*Card) {
	l.cards = append(cards, l.cards...)
}

// PutOnBottom places the specified cards on the bottom of the library.
func (l *Library) PutOnBottom(cards ...*Card) {
	l.cards = append(l.cards, cards...)
}

// Size returns the number of cards in the library.
func (l *Library) Size() int {
	return len(l.cards)
}
