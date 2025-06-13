package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
	"fmt"
)

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
	return p
}

func (p Player) WithManaPool(manaPool mana.Pool) Player {
	p.manaPool = manaPool
	return p
}

func (p Player) WithEmptyManaPool() Player {
	p.manaPool = mana.NewManaPool()
	return p
}

func (p Player) WithAddCardToZone(card gob.Card, zone mtg.Zone) (Player, bool) {
	switch zone {
	case mtg.ZoneExile:
		exile := p.WithExile(p.exile.Add(card))
		p = exile
	case mtg.ZoneGraveyard:
		graveyard := p.WithGraveyard(p.graveyard.Add(card))
		p = graveyard
	case mtg.ZoneHand:
		hand := p.WithHand(p.hand.Add(card))
		p = hand
	case mtg.ZoneLibrary:
		library := p.WithLibrary(p.library.Add(card))
		p = library
	default:
		return p, false // No change for unsupported zones
	}
	return p, true
}

func (p Player) WithAddCardToTopOfZone(card gob.Card, zone mtg.Zone) (Player, bool) {
	switch zone {
	case mtg.ZoneExile:
		exile := p.WithExile(p.exile.AddTop(card))
		p = exile
	case mtg.ZoneGraveyard:
		graveyard := p.WithGraveyard(p.graveyard.AddTop(card))
		p = graveyard
	case mtg.ZoneHand:
		hand := p.WithHand(p.hand.AddTop(card))
		p = hand
	case mtg.ZoneLibrary:
		library := p.WithLibrary(p.library.AddTop(card))
		p = library
	default:
		return p, false // No change for unsupported zones
	}
	return p, true
}

// TODO: Should this be WithMoveCardToZone?
func (p Player) WithDiscardCard(cardID string) (Player, error) {
	card, newHand, ok := p.hand.Take(cardID)
	if !ok {
		return p, fmt.Errorf("card %q not found", cardID)
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

func (p Player) WithGainLife(amount int) Player {
	p.life += amount
	return p
}

func (p Player) WithLoseLife(amount int) Player {
	p.life -= amount
	return p
}

func (p Player) WithClearLandPlayedThisTurn() Player {
	p.landPlayedThisTurn = false
	return p
}

func (p Player) WithRevealed(revealed Revealed) Player {
	p.revealed = revealed
	return p
}

func (p Player) WithClearRevealed() Player {
	p.revealed = NewRevealed()
	return p
}

func (p Player) WithSpellCastThisTurn() Player {
	p.spellsCastThisTurn++
	return p
}

func (p Player) WithClearSpellsCastsThisTurn() Player {
	p.spellsCastThisTurn = 0
	return p
}
