package game

import "errors"

type Revealed struct {
	cards []*Card
}

func NewRevealed() *Revealed {
	return &Revealed{
		cards: []*Card{},
	}
}

func (r *Revealed) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return errors.New("object is not a card")
	}
	r.cards = append(r.cards, card)
	return nil
}

// This probably makes more sense as a method of Player
func (r *Revealed) AvailableActivatedAbilities(*GameState, PlayerAgent) []*ActivatedAbility {
	return nil
}

// This probably makes more sense as a method of Player
func (r *Revealed) AvailableToPlay(*GameState, PlayerAgent) []GameObject {
	return nil
}

func (r *Revealed) Clear() {
	r.cards = []*Card{}
}

func (r *Revealed) Find(id string) (GameObject, error) {
	return nil, nil
}

func (r *Revealed) Get(id string) (GameObject, error) {
	return nil, nil
}

func (r *Revealed) GetAll() []GameObject {
	var all []GameObject
	for _, card := range r.cards {
		all = append(all, card)
	}
	return all
}

func (r *Revealed) Remove(id string) error {
	return nil
}

func (r *Revealed) Take(id string) (GameObject, error) {
	return nil, nil
}

func (r *Revealed) Size() int {
	return 0
}

func (r *Revealed) ZoneType() string {
	return ZoneRevealed
}
