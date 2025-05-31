package zone

import (
	"deckronomicon/packages/game/card"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Graveyard struct {
	cards []*card.Card
}

// NewGraveyard creates a new Graveyard instance.
func NewGraveyard() *Graveyard {
	return &Graveyard{
		cards: []*card.Card{},
	}
}

func (g *Graveyard) Add(card *card.Card) {
	g.cards = append(g.cards, card)
}

func (g *Graveyard) Get(id string) (*card.Card, error) {
	for _, card := range g.cards {
		if card.ID() == id {
			return card, nil
		}
	}
	return nil, fmt.Errorf("card witg ID %s not found", id)
}

func (g *Graveyard) GetAll() []*card.Card {
	return g.cards
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

func (g *Graveyard) Take(id string) (*card.Card, error) {
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

func (g *Graveyard) ZoneType() mtg.Zone {
	return mtg.ZoneGraveyard
}

/*
// TODO: Move to player
// AvailableActivatedAbilities returns all activated abilities of the cards in
// the graveyard.
func (g *Graveyard) AvailableActivatedAbilities(state *GameState, player *Player) []game.Object {
	var objects []game.Object
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
			objects = append(objects, ability)
		}
	}
	return objects
}
*/
/*
// AvailableToPlay returns a list of cards that can be played from the
// graveyard. Static abilities such as Flashback and Unearth are not
// implemented yet so this wil return nil.
func (g *Graveyard) AvailableToPlay(state *GameState, player *Player) []game.Object {
	var found []game.Object
	for _, card := range g.cards {
		fmt.Println("Checking card:", card.ID())
		for _, ability := range card.StaticAbilities() {
			fmt.Println("Checking ability:", card.ID())
			if ability.ID == game.AbilityKeywordFlashback {
				fmt.Println("Found Flashback ability on card:", card.ID())
				found = append(found, card)
			}
		}
	}
	return found
}
*/
