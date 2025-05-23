package game

import (
	"fmt"
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

// AvailableActivatedAbilities returns a list of activated abilities that can
// be activated from the library. This exits to satisfy the Zone interface.
// Cards in the library cannot have a activated ability.
func (l *Library) AvailableActivatedAbilities(*GameState) []*ActivatedAbility {
	return nil
}

// AvailableToPlay returns a list of cards that can be played from the
// library. This exists to satisfy the Zone interface. Cards in the library
// generally cannot be played. (Some cards enable this, but they are
// exceptions and not currently implemented.)
func (l *Library) AvailableToPlay(*GameState) []GameObject {
	return nil
}

func (l *Library) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return fmt.Errorf("object is not a card")
	}
	l.cards = append(l.cards, card)
	return nil
}

func (l *Library) AddTop(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return fmt.Errorf("object is not a card")
	}
	l.cards = append([]*Card{card}, l.cards...)
	return nil
}

func (l *Library) Find(id string) (GameObject, error) {
	for _, card := range l.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

// FindByName finds the first card in the hand by name.
func (l *Library) FindByName(name string) (GameObject, error) {
	for _, card := range l.cards {
		if card.Name() == name {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with name %s not found", name)
}

// FindAllBySubtype finds all cards in the library by subtype.
func (l *Library) FindAllBySubtype(subtype Subtype) []GameObject {
	var foundCards []GameObject
	for _, card := range l.cards {
		if card.HasSubtype(subtype) {
			foundCards = append(foundCards, card)
		}
	}
	return foundCards
}

func (l *Library) Get(id string) (GameObject, error) {
	for _, card := range l.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (l *Library) GetAll() []GameObject {
	var all []GameObject
	for _, card := range l.cards {
		all = append(all, card)
	}
	return all
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

func (l *Library) Take(id string) (GameObject, error) {
	for i, card := range l.cards {
		if card.ID() == id {
			l.cards = append(l.cards[:i], l.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

// TakeCards removes the top N cards from the library and returns them.
func (l *Library) TakeTop() (GameObject, error) {
	if len(l.cards) == 0 {
		return nil, ErrLibraryEmpty
	}
	l.cards = l.cards[1:]
	return l.cards[0], nil
}

func (l *Library) Size() int {
	return len(l.cards)
}

func (l *Library) ZoneType() string {
	return ZoneLibrary
}

// Library Specific Methods

// Peek returns the top N cards without modifying the library.
func (l *Library) Peek() *Card {
	if len(l.cards) == 0 {
		return nil
	}
	return l.cards[0]
}

// Shuffle randomly shuffles the cards in the library.
func (l *Library) Shuffle() {
	rand.Shuffle(len(l.cards), func(i, j int) {
		l.cards[i], l.cards[j] = l.cards[j], l.cards[i]
	})
}
