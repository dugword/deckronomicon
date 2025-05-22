package game

import "fmt"

type Graveyard struct {
	cards []*Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() *Graveyard {
	return &Graveyard{
		cards: []*Card{},
	}
}

// TODO Remove access to this
func (g *Graveyard) Cards() []*Card {
	return g.cards
}

func (g *Graveyard) Add(object GameObject) error {
	card, ok := object.(*Card)
	if !ok {
		return fmt.Errorf("object is not a card")
	}
	g.cards = append(g.cards, card)
	return nil
}

// AvailableActivatedAbilities returns all activated abilities of the cards in
// the graveyard.
func (g *Graveyard) AvailableActivatedAbilities(state *GameState) []*ActivatedAbility {
	abilities := []*ActivatedAbility{}
	for _, card := range g.cards {
		for _, ability := range card.ActivatedAbilities() {
			if !ability.Cost.CanPay(state) {
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
func (g *Graveyard) AvailableToPlay(state *GameState) []GameObject {
	return nil
}

func (g *Graveyard) Find(id string) (GameObject, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

// FindByName finds the first card in the graveyard by name.
func (g *Graveyard) FindByName(name string) (GameObject, error) {
	for _, card := range g.cards {
		if card.Name() == name {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card with name %s not found", name)
}

// FindAllBySubtype finds all cards in the graveyard by subtype.
func (g *Graveyard) FindAllBySubtype(subtype Subtype) []GameObject {
	cards := []GameObject{}
	for _, card := range g.cards {
		if card.HasSubtype(subtype) {
			cards = append(cards, card)
		}
	}
	return cards
}

func (g *Graveyard) Get(id string) (GameObject, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
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
