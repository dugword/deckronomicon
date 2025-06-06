package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

// Battlefield represents the battlefield during an active game.
type Battlefield struct {
	permanents []gob.Permanent
}

// NewBattlefield creates a new Battlefield instance.
func NewBattlefield() Battlefield {
	battlefield := Battlefield{
		permanents: []gob.Permanent{},
	}
	return battlefield
}

func (b Battlefield) Add(permanents ...gob.Permanent) Battlefield {
	newPermanents := append(b.permanents[:], permanents...)
	return Battlefield{
		permanents: newPermanents,
	}
}

func (b Battlefield) Get(id string) (gob.Permanent, error) {
	for _, permanent := range b.permanents {
		if permanent.ID() == id {
			return permanent, nil
		}
	}
	return gob.Permanent{}, fmt.Errorf("permanent with ID %s not found", id)
}

func (b Battlefield) GetAll() []gob.Permanent {
	return b.permanents
}

func (b Battlefield) Name() string {
	return string(mtg.ZoneBattlefield)
}

func (b Battlefield) Remove(id string) (Battlefield, error) {
	for i, permanent := range b.permanents {
		if permanent.ID() == id {
			b.permanents = append(b.permanents[:i], b.permanents[i+1:]...)
			return b, nil
		}
	}
	return b, fmt.Errorf("permanent with ID %s not found", id)
}

func (b Battlefield) Size() int {
	return len(b.permanents)
}

func (b Battlefield) Take(id string) (gob.Permanent, Battlefield, error) {
	for i, permanent := range b.permanents {
		if permanent.ID() == id {
			b.permanents = append(b.permanents[:i], b.permanents[i+1:]...)
			return permanent, b, nil
		}
	}
	return gob.Permanent{}, b, fmt.Errorf("permanent with ID %s not found", id)
}

func (b Battlefield) UntapAll(playerID string) Battlefield {
	var battlefield Battlefield
	for _, p := range b.permanents {
		if p.Controller() == playerID {
			p.Untap()
		}
		battlefield.Add(p)
	}
	return battlefield
}

func (b Battlefield) Zone() mtg.Zone {
	return mtg.ZoneBattlefield
}
