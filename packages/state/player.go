package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
)

type Player struct {
	exile       Exile
	graveyard   Graveyard
	hand        Hand
	id          string
	landDrop    bool
	library     Library
	life        int
	maxHandSize int // TODO make this configurable
	turn        int
}

func (p Player) Life() int {
	return p.life
}

func (p Player) Turn() int {
	return p.turn
}

func NewPlayer(id string, deckList []gob.Card) Player {
	return Player{
		exile:       NewExile(),
		graveyard:   NewGraveyard(),
		hand:        NewHand(),
		id:          id,
		library:     NewLibrary(deckList),
		life:        20, // TODO make this configurable
		maxHandSize: 7,
	}
}

func (p Player) TakeCardFromHand(
	cardID string,
) (gob.Card, Player, error) {
	card, newHand, ok := p.hand.Take(cardID)
	if !ok {
		return card, p, fmt.Errorf("card %s not found in hand", cardID)
	}
	player := p.WithHand(newHand)
	return card, player, nil
}

func (p Player) MaxHandSize() int {
	return p.maxHandSize
}

func (p Player) WithShuffleDeck(
	deckShuffler func([]gob.Card) []gob.Card,
) Player {
	// Shuffle the library
	cards := deckShuffler(p.library.GetAll())
	newLibrary := NewLibrary(cards)
	p.library = newLibrary
	return p
}

func (p Player) WithNextTurn() Player {
	p.turn++
	p.landDrop = false
	return p
}

func (p Player) WithDiscardCard(cardID string) (Player, error) {
	card, newHand, ok := p.hand.Take(cardID)
	if !ok {
		return p, fmt.Errorf("card %s not found in hand", cardID)
	}
	newGraveyard := p.graveyard.Append(card)
	player := p.WithHand(newHand).WithGraveyard(newGraveyard)
	return player, nil
}

func (p Player) WithDrawCard() (Player, gob.Card, error) {
	card, library, ok := p.library.Shift()
	if !ok {
		return p, gob.Card{}, mtg.ErrLibraryEmpty
	}
	hand := p.hand.Append(card)
	player := p.WithLibrary(library).WithHand(hand)
	return player, card, nil
}

func (p Player) WithLibrary(library Library) Player {
	p.library = library
	return p
}

func (p Player) WithHand(hand Hand) Player {
	p.hand = hand
	return p
}

func (p Player) WithExile(exile Exile) Player {
	p.exile = exile
	return p
}

func (p Player) WithGraveyard(graveyard Graveyard) Player {
	p.graveyard = graveyard
	return p
}

func (p Player) ID() string {
	return p.id
}

// TODO: Do I still need this view abstraction after the refactor?
// Or since everything is read only and immutable, can I just return the
// struct directly?
func (p Player) Library() query.View {
	return query.NewView(
		string(mtg.ZoneLibrary),
		p.library.GetAll(),
	)
}

func (p Player) Hand() Hand {
	return p.hand
}

/*
func (p Player) Hand() query.View {
	return query.NewView(
		string(mtg.ZoneHand),
		p.hand.GetAll(),
	)
}
*/

func (p Player) Exile() query.View {
	return query.NewView(
		string(mtg.ZoneExile),
		p.exile.GetAll(),
	)
}

func (p Player) Graveyard() query.View {
	return query.NewView(
		string(mtg.ZoneGraveyard),
		p.graveyard.GetAll(),
	)
}
