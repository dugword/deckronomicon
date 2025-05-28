package game

import "errors"

type Exile struct {
	cards []*Card
}

func NewExile() *Exile {
	return &Exile{
		cards: []*Card{},
	}
}

func (e *Exile) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return errors.New("object is not a card")
	}
	e.cards = append(e.cards, card)
	return nil
}

// This probably makes more sense as a method of Player
func (e *Exile) AvailableActivatedAbilities(*GameState, PlayerAgent) []GameObject {
	return nil
}

// This probably makes more sense as a method of Player
func (e *Exile) AvailableToPlay(*GameState, PlayerAgent) []GameObject {
	return nil
}
func (e *Exile) Find(id string) (GameObject, error) {
	for _, card := range e.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}
func (e *Exile) Get(id string) (GameObject, error) {
	for _, card := range e.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}
func (e *Exile) GetAll() []GameObject {
	var objects []GameObject
	for _, card := range e.cards {
		objects = append(objects, card)
	}
	return objects
}
func (e *Exile) Remove(id string) error {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return nil
		}
	}
	return errors.New("card not found in exile")
}
func (e *Exile) Take(id string) (GameObject, error) {
	for i, card := range e.cards {
		if card.ID() == id {
			e.cards = append(e.cards[:i], e.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, errors.New("card not found in exile")
}
func (e *Exile) Size() int {
	return len(e.cards)
}
func (e *Exile) ZoneType() string {
	return ZoneExile
}
