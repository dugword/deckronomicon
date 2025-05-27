package game

import (
	"errors"
	"fmt"
)

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
		// TODO: Move all errors.New to the errors file
		return errors.New("object is not a card")
	}
	h.cards = append(h.cards, card)
	return nil
}

func (h *Hand) AvailableActivatedAbilities(state *GameState, player *Player) []*ActivatedAbility {
	var abilities []*ActivatedAbility
	for _, card := range h.cards {
		for _, ability := range card.ActivatedAbilities() {
			if !ability.CanPlay(state) {
				continue
			}
			if !ability.Cost.CanPay(state, player) {
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

func (h *Hand) AvailableToPlay(state *GameState, player *Player) []GameObject {
	var available []GameObject
	for _, card := range h.cards {
		if state.CanCastSorcery() {
			if card.IsLand() {
				if player.LandDrop == false {
					available = append(available, card)
					continue
				}
			} else if card.ManaCost().CanPay(state, player) {
				available = append(available, card)
				continue
			}
		}
		if card.ManaCost().CanPay(state, player) && card.HasCardType(CardTypeInstant) {
			available = append(available, card)
			continue
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
