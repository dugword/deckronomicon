package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query"
	"fmt"
)

type Player struct {
	exile              Exile
	graveyard          Graveyard
	hand               Hand
	id                 string
	landPlayedThisTurn bool
	library            Library
	life               int
	manaPool           mana.Pool
	maxHandSize        int // TODO make this configurable
	turn               int
}

func (p Player) Life() int {
	return p.life
}

func (p Player) Turn() int {
	return p.turn
}

func NewPlayer(id string, life int) Player {
	return Player{
		exile:       NewExile(),
		graveyard:   NewGraveyard(),
		hand:        NewHand(),
		id:          id,
		library:     NewLibrary([]gob.Card{}), // Start with an empty library
		life:        life,
		manaPool:    mana.NewManaPool(),
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

func (p Player) WithAddMana(manaType mana.ManaType, amount int) Player {
	p.manaPool = p.manaPool.WithAddedMana(manaType, amount)
	return p
}

func (p Player) WithShuffleDeck(
	deckShuffler func([]gob.Card) []gob.Card,
) Player {
	// Shuffle the library
	cards := deckShuffler(p.library.GetAll())
	p.library.cards = cards
	return p
}

func (p Player) WithNextTurn() Player {
	p.turn++
	p.landPlayedThisTurn = false
	return p
}

func (p Player) WithEmptyManaPool() Player {
	p.manaPool = mana.NewManaPool()
	return p
}

func (p Player) WithDiscardCard(cardID string) (Player, error) {
	card, newHand, ok := p.hand.Take(cardID)
	if !ok {
		return p, fmt.Errorf("card %s not found in hand", cardID)
	}
	newGraveyard := p.graveyard.Add(card)
	player := p.WithHand(newHand).WithGraveyard(newGraveyard)
	return player, nil
}

func (p Player) WithDrawCard() (Player, gob.Card, error) {
	card, library, ok := p.library.TakeTop()
	if !ok {
		return p, gob.Card{}, mtg.ErrLibraryEmpty
	}
	hand := p.hand.Add(card)
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

func (p Player) WithLandPlayedThisTurn() Player {
	p.landPlayedThisTurn = true
	return p
}

func (p Player) LandPlayedThisTurn() bool {
	return p.landPlayedThisTurn
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

func (p Player) ManaPool() mana.Pool {
	return p.manaPool
}

/*
func (p Player) Hand() query.View {
	return query.NewView(
		string(mtg.ZoneHand),
		p.hand.GetAll(),
	)
}
*/

func (p Player) Exile() Exile {
	return p.exile
}

func (p Player) Graveyard() Graveyard {
	return p.graveyard
}
