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
	spellsCastThisTurn int
	turn               int
	revealed           Revealed
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

func (p Player) TakeCardFromZone(
	cardID string,
	zone mtg.Zone,
) (gob.Card, Player, bool) {
	var card gob.Card
	switch zone {
	case mtg.ZoneExile:
		c, exile, ok := p.exile.Take(cardID)
		if !ok {
			return card, p, false
		}
		p = p.WithExile(exile)
		card = c
	case mtg.ZoneGraveyard:
		c, graveyard, ok := p.graveyard.Take(cardID)
		if !ok {
			return card, p, false
		}
		p = p.WithGraveyard(graveyard)
		card = c
	case mtg.ZoneHand:
		c, hand, ok := p.hand.Take(cardID)
		if !ok {
			return card, p, false
		}
		p = p.WithHand(hand)
		card = c
	case mtg.ZoneLibrary:
		c, library, ok := p.library.Take(cardID)
		if !ok {
			return card, p, false
		}
		p = p.WithLibrary(library)
		card = c
	default:
		return card, p, false
	}
	return card, p, true
}

func (p Player) MaxHandSize() int {
	return p.maxHandSize
}

func (p Player) SpellsCastThisTurn() int {
	return p.spellsCastThisTurn
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
func (p Player) Library() Library {
	return p.library
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

func (p Player) Revealed() Revealed {
	return p.revealed
}

func (p Player) GetZone(zone mtg.Zone) (query.View, bool) {
	switch zone {
	case mtg.ZoneExile:
		return query.NewView(string(mtg.ZoneExile), p.exile.GetAll()), true
	case mtg.ZoneGraveyard:
		return query.NewView(string(mtg.ZoneGraveyard), p.graveyard.GetAll()), true
	case mtg.ZoneHand:
		return query.NewView(string(mtg.ZoneHand), p.hand.GetAll()), true
	case mtg.ZoneLibrary:
		return query.NewView(string(mtg.ZoneLibrary), p.library.GetAll()), true
	default:
		return nil, false
	}
}
