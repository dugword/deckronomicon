package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

type Player struct {
	exile              *Exile
	graveyard          *Graveyard
	hand               *Hand
	id                 string
	landPlayedThisTurn bool
	library            *Library
	life               int
	manaPool           mana.Pool
	maxHandSize        int
	spellsCastThisTurn int
	turn               int
	revealed           *Revealed
}

func NewPlayer(id string, life int) *Player {
	return &Player{
		exile:       NewExile(),
		graveyard:   NewGraveyard(),
		hand:        NewHand(),
		id:          id,
		library:     NewLibrary([]*gob.Card{}),
		life:        life,
		manaPool:    mana.Pool{},
		maxHandSize: 7,
	}
}

func (p *Player) Name() string {
	return p.id
}

func (p *Player) Life() int {
	return p.life
}

func (p *Player) Turn() int {
	return p.turn
}

func (p *Player) GetCardFromZone(
	cardID string,
	zone mtg.Zone,
) (*gob.Card, bool) {
	switch zone {
	case mtg.ZoneExile:
		return p.exile.Get(cardID)
	case mtg.ZoneGraveyard:
		return p.graveyard.Get(cardID)
	case mtg.ZoneHand:
		return p.hand.Get(cardID)
	case mtg.ZoneLibrary:
		return p.library.Get(cardID)
	default:
		return nil, false
	}
}

func (p *Player) TakeCardFromZone(
	cardID string,
	zone mtg.Zone,
) (*gob.Card, *Player, bool) {
	switch zone {
	case mtg.ZoneExile:
		c, exile, ok := p.exile.Take(cardID)
		if !ok {
			return nil, nil, false
		}
		return c, p.WithExile(exile), true
	case mtg.ZoneGraveyard:
		c, graveyard, ok := p.graveyard.Take(cardID)
		if !ok {
			return nil, nil, false
		}
		return c, p.WithGraveyard(graveyard), true
	case mtg.ZoneHand:
		c, hand, ok := p.hand.Take(cardID)
		if !ok {
			return nil, nil, false
		}
		return c, p.WithHand(hand), true
	case mtg.ZoneLibrary:
		c, library, ok := p.library.Take(cardID)
		if !ok {
			return nil, nil, false
		}
		return c, p.WithLibrary(library), true
	default:
		return nil, nil, false
	}
}

func (p *Player) MaxHandSize() int {
	return p.maxHandSize
}

func (p *Player) SpellsCastThisTurn() int {
	return p.spellsCastThisTurn
}

func (p *Player) LandPlayedThisTurn() bool {
	return p.landPlayedThisTurn
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Library() *Library {
	return p.library
}

func (p *Player) Hand() *Hand {
	return p.hand
}

func (p *Player) ManaPool() mana.Pool {
	return p.manaPool.Copy()
}

func (p *Player) Exile() *Exile {
	return p.exile
}

func (p *Player) Graveyard() *Graveyard {
	return p.graveyard
}

func (p *Player) Revealed() *Revealed {
	return p.revealed
}

func (p *Player) GetCardsInZone(zone mtg.Zone) ([]*gob.Card, bool) {
	switch zone {
	case mtg.ZoneExile:
		return p.exile.GetAll(), true
	case mtg.ZoneGraveyard:
		return p.graveyard.GetAll(), true
	case mtg.ZoneHand:
		return p.hand.GetAll(), true
	case mtg.ZoneLibrary:
		return p.library.GetAll(), true
	default:
		return nil, false
	}
}

func (p *Player) ZoneContains(
	zone mtg.Zone,
	predicate query.Predicate,
) bool {
	switch zone {
	case mtg.ZoneExile:
		return p.exile.Contains(predicate)
	case mtg.ZoneGraveyard:
		return p.graveyard.Contains(predicate)
	case mtg.ZoneHand:
		return p.hand.Contains(predicate)
	case mtg.ZoneLibrary:
		return p.library.Contains(predicate)
	default:
		return false
	}
}
