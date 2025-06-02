package zone

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/permanent"
	"fmt"
)

// Battlefield represents the battlefield during an active game.
type Battlefield struct {
	permanents []*permanent.Permanent
}

// NewBattlefield creates a new Battlefield instance.
func NewBattlefield() *Battlefield {
	battlefield := Battlefield{
		permanents: []*permanent.Permanent{},
	}
	return &battlefield
}

func (b *Battlefield) Add(permanent *permanent.Permanent) {
	b.permanents = append(b.permanents, permanent)
}

func (b *Battlefield) Get(id string) (*permanent.Permanent, error) {
	for _, permanent := range b.permanents {
		if permanent.ID() == id {
			return permanent, nil
		}
	}
	return nil, fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) GetAll() []*permanent.Permanent {
	return b.permanents
}

func (b *Battlefield) Name() string {
	return string(mtg.ZoneBattlefield)
}

func (b *Battlefield) Remove(id string) error {
	for i, permanent := range b.permanents {
		if permanent.ID() == id {
			b.permanents = append(b.permanents[:i], b.permanents[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) Size() int {
	return len(b.permanents)
}

func (b *Battlefield) Take(id string) (*permanent.Permanent, error) {
	for i, permanent := range b.permanents {
		if permanent.ID() == id {
			b.permanents = append(b.permanents[:i], b.permanents[i+1:]...)
			return permanent, nil
		}
	}
	return nil, fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) Zone() mtg.Zone {
	return mtg.ZoneBattlefield
}
