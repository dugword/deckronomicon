package game

import "fmt"

// Hand represents a player's hand of cards.
type Hand struct {
	cards []*Card
}

// NewHand creates a new Hand instance.
func NewHand() *Hand {
	return &Hand{
		cards: []*Card{},
	}
}

func (h *Hand) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return fmt.Errorf("object is not a card")
	}
	h.cards = append(h.cards, card)
	return nil
}

func (h *Hand) AvailableActivatedAbilities(state *GameState) []*ActivatedAbility {
	var abilities []*ActivatedAbility
	for _, card := range h.cards {
		for _, ability := range card.ActivatedAbilities() {
			if !ability.Cost.CanPay(state) {
				continue
			}
			if ability.Zone != ZoneHand {
				continue

			}
			abilities = append(abilities, ability)
		}
	}
	return abilities
}

func (h *Hand) AvailableToPlay(state *GameState) []GameObject {
	var available []GameObject
	for _, card := range h.cards {
		if (card.IsLand() && state.LandDrop == false) || card.ManaCost().CanPay(state) {
			available = append(available, card)
		}
	}
	return available
}

func (h *Hand) Find(id string) (GameObject, error) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

// FindByName finds the first card in the hand by name.
func (h *Hand) FindByName(name string) (GameObject, error) {
	for _, card := range h.cards {
		if card.Name() == name {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with name %s not found", name)
}

// FindAllBySubtype finds all cards in the hand by subtype.
func (h *Hand) FindAllBySubtype(subtype Subtype) []GameObject {
	var found []GameObject
	for _, card := range h.cards {
		if card.HasSubtype(subtype) {
			found = append(found, card)
		}
	}
	return found
}

func (h *Hand) Get(id string) (GameObject, error) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) GetAll() []GameObject {
	var all []GameObject
	for _, card := range h.cards {
		all = append(all, card)
	}
	return all
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

func (h *Hand) Take(id string) (GameObject, error) {
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

func (h *Hand) ZoneType() string {
	return ZoneHand
}

// Cards returns the cards in the hand.
// TODO: for display remove later - why? Probably should directly manipulate
// this
func (h *Hand) Cards() []*Card {
	return h.cards
}
