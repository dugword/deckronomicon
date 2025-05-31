package zone

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

// Hand represents a player's hand of cards.
type Hand struct {
	cards []*card.Card
}

// NewHand creates a new Hand instance.
func NewHand() *Hand {
	return &Hand{
		cards: []*card.Card{},
	}
}

func (h *Hand) Add(card *card.Card) {
	h.cards = append(h.cards, card)
}

func (h *Hand) Find(id string) (*card.Card, error) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) Get(id string) (*card.Card, error) {
	for _, card := range h.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with ID %s not found", id)
}

func (h *Hand) GetAll() []*card.Card {
	return h.cards
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

func (h *Hand) Take(id string) (*card.Card, error) {
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

/*
func (h *Hand) AvailableActivatedAbilities(state *GameState, player *Player) []game.Object {
	var objects []game.Object
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
			objects = append(objects, ability)
		}
	}
	return objects
}

func (h *Hand) AvailableToPlay(state *GameState, player *Player) []game.Object {
	var available []game.Object
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
		if card.ManaCost().CanPay(state, player) && card.HasCardType(game.CardTypeInstant) {
			available = append(available, card)
			continue
		}
	}
	return available
}
*/
