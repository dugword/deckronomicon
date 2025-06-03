package engine

import "deckronomicon/packages/game/gob"

type Player struct {
	id        string
	passed    bool
	name      string
	life      int
	hand      []*gob.Card
	library   []*gob.Card
	graveyard []*gob.Card
	exile     []*gob.Card
}

func NewPlayer(id string, deckList []*gob.Card) *Player {
	return &Player{
		id:      id,
		library: deckList,
		life:    20, // TODO make this configurable
	}
}
