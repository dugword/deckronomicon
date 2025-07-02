package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

func (p *Player) WithAddMana(color mana.Color, amount int) *Player {
	newPlayer := *p
	manaPool := newPlayer.manaPool.WithAddMana(amount, color)
	newPlayer.manaPool = manaPool
	return &newPlayer
}

func (p *Player) WithNextTurn() *Player {
	newPlayer := *p
	newPlayer.turn++
	return &newPlayer
}

func (p *Player) WithManaPool(manaPool mana.Pool) *Player {
	newPlayer := *p
	newPlayer.manaPool = manaPool
	return &newPlayer
}

func (p *Player) WithEmptyManaPool() *Player {
	newPlayer := *p
	newPlayer.manaPool = mana.Pool{}
	return &newPlayer
}

func (p *Player) WithAddCardToZone(card *gob.Card, zone mtg.Zone) (*Player, bool) {
	switch zone {
	case mtg.ZoneExile:
		return p.WithExile(p.exile.Add(card)), true
	case mtg.ZoneGraveyard:
		return p.WithGraveyard(p.graveyard.Add(card)), true
	case mtg.ZoneHand:
		return p.WithHand(p.hand.Add(card)), true
	case mtg.ZoneLibrary:
		return p.WithLibrary(p.library.Add(card)), true
	default:
		return nil, false // No change for unsupported zones
	}
}

func (p *Player) WithAddCardToTopOfZone(card *gob.Card, zone mtg.Zone) (*Player, bool) {
	switch zone {
	case mtg.ZoneExile:
		return p.WithExile(p.exile.AddTop(card)), true
	case mtg.ZoneGraveyard:
		return p.WithGraveyard(p.graveyard.AddTop(card)), true
	case mtg.ZoneHand:
		return p.WithHand(p.hand.AddTop(card)), true
	case mtg.ZoneLibrary:
		return p.WithLibrary(p.library.AddTop(card)), true
	default:
		return nil, false // No change for unsupported zones
	}
}

func (p *Player) WithDiscardCard(cardID string) (*Player, error) {
	card, newHand, ok := p.hand.Take(cardID)
	if !ok {
		return nil, fmt.Errorf("card %q not found", cardID)
	}
	newGraveyard := p.graveyard.Add(card)
	newPlayer := p.WithHand(newHand).WithGraveyard(newGraveyard)
	return newPlayer, nil
}

func (p *Player) WithDrawCard() (*Player, *gob.Card, error) {
	card, library, ok := p.library.TakeTop()
	if !ok {
		return nil, nil, mtg.ErrLibraryEmpty
	}
	hand := p.hand.Add(card)
	player := p.WithLibrary(library).WithHand(hand)
	return player, card, nil
}

func (p *Player) WithLibrary(library *Library) *Player {
	newPlayer := *p
	newPlayer.library = library
	return &newPlayer
}

func (p *Player) WithHand(hand *Hand) *Player {
	newPlayer := *p
	newPlayer.hand = hand
	return &newPlayer
}

func (p *Player) WithExile(exile *Exile) *Player {
	newPlayer := *p
	newPlayer.exile = exile
	return &newPlayer
}

func (p *Player) WithGraveyard(graveyard *Graveyard) *Player {
	newPlayer := *p
	newPlayer.graveyard = graveyard
	return &newPlayer
}

func (p *Player) WithLandPlayedThisTurn() *Player {
	newPlayer := *p
	newPlayer.landPlayedThisTurn = true
	return &newPlayer
}

func (p *Player) WithGainLife(amount int) *Player {
	newPlayer := *p
	newPlayer.life += amount
	return &newPlayer
}

func (p *Player) WithLoseLife(amount int) *Player {
	newPlayer := *p
	newPlayer.life -= amount
	return &newPlayer
}

func (p *Player) WithClearLandPlayedThisTurn() *Player {
	newPlayer := *p
	newPlayer.landPlayedThisTurn = false
	return &newPlayer
}

func (p *Player) WithRevealed(revealed *Revealed) *Player {
	newPlayer := *p
	newPlayer.revealed = revealed
	return &newPlayer
}

func (p *Player) WithClearRevealed() *Player {
	newPlayer := *p
	newPlayer.revealed = NewRevealed()
	return &newPlayer
}

func (p *Player) WithSpellCastThisTurn() *Player {
	newPlayer := *p
	newPlayer.spellsCastThisTurn++
	return &newPlayer
}

func (p *Player) WithShuffledLibrary(shuffledCardsIDs []string) (*Player, error) {
	cards := p.library.GetAll()
	if len(cards) != len(shuffledCardsIDs) {
		return nil, fmt.Errorf(
			"shuffledCardIDs %d does not match library size %d",
			len(shuffledCardsIDs),
			len(cards),
		)
	}
	cardIDMap := map[string]*gob.Card{}
	for _, card := range cards {
		cardIDMap[card.ID()] = card
	}
	var redordered []*gob.Card
	for _, id := range shuffledCardsIDs {
		if card, ok := cardIDMap[id]; ok {
			redordered = append(redordered, card)
		} else {
			return nil, fmt.Errorf("card %q not found in library", id)
		}
	}
	return p.WithLibrary(&Library{cards: redordered}), nil
}

func (p *Player) WithClearSpellsCastsThisTurn() *Player {
	newPlayer := *p
	newPlayer.spellsCastThisTurn = 0
	return &newPlayer
}
