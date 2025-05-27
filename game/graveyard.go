package game

import (
	"errors"
	"fmt"
)

type Graveyard struct {
	cards []*Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() *Graveyard {
	return &Graveyard{
		cards: []*Card{},
	}
}

func (g *Graveyard) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return errors.New("object is not a card")
	}
	g.cards = append(g.cards, card)
	return nil
}

// AvailableActivatedAbilities returns all activated abilities of the cards in
// the graveyard.
func (g *Graveyard) AvailableActivatedAbilities(state *GameState, player *Player) []*ActivatedAbility {
	abilities := []*ActivatedAbility{}
	for _, card := range g.cards {
		for _, ability := range card.ActivatedAbilities() {
			if !ability.CanPlay(state) {
				continue
			}
			if !ability.Cost.CanPay(state, player) {
				continue
			}
			if ability.Zone != ZoneGraveyard {
				continue
			}
			abilities = append(abilities, ability)
		}
	}
	return abilities
}

// AvailableToPlay returns a list of cards that can be played from the
// graveyard. Static abilities such as Flashback and Unearth are not
// implemented yet so this wil return nil.
func (g *Graveyard) AvailableToPlay(state *GameState, player *Player) []GameObject {
	var found []GameObject
	for _, card := range g.cards {
		fmt.Println("Checking card:", card.ID())
		for _, ability := range card.StaticAbilities() {
			fmt.Println("Checking ability:", card.ID())
			if ability.ID == AbilityKeywordFlashback {
				fmt.Println("Found Flashback ability on card:", card.ID())
				found = append(found, card)
			}
		}
	}
	return found
}

func (g *Graveyard) Find(id string) (GameObject, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) Get(id string) (GameObject, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) GetAll() []GameObject {
	var all []GameObject
	for _, card := range g.cards {
		all = append(all, card)
	}
	return all
}

func (g *Graveyard) Remove(id string) error {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) Take(id string) (GameObject, error) {
	for i, card := range g.cards {
		if card.ID() == id {
			g.cards = append(g.cards[:i], g.cards[i+1:]...)
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) Size() int {
	return len(g.cards)
}

func (g *Graveyard) ZoneType() string {
	return ZoneGraveyard
}
