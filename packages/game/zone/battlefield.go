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

// Battlefield Specific Methods

/*
// TODO Think if this makes sense or if it should run in the untap step,
// collect all things and then remove summoning sickness and untap them

// RemoveSummoningSickness removes summoning sickness from all permanents
func (b *Battlefield) RemoveSummoningSickness() {
	for _, permanent := range b.permanents {
		permanent.RemoveSummoningSickness()
	}
}

// UntapPermanents untaps all permanents on the battlefield.
// TODO: In the future this will support managing only untapping permanents
// that don't have effects/stun counters preventing them from being untapped.
// Or maybe that would be handled by the .Untap method?
func (b *Battlefield) UntapPermanents() {
	for _, permanent := range b.permanents {
		permanent.Untap()
	}
}
*/

// TODO: move to player
/*
func (b *Battlefield) AvailableActivatedAbilities(state state, player player) []activated.Ability {
	var abilities []activated.Ability
	for _, permanent := range b.permanents {
		for _, ability := range permanent.ActivatedAbilities() {
			if !ability.CanActivate(state) {
				continue
			}
			if ability.Cost.CanPay(state, player) {
				objects = append(objects, ability)
			}
		}
	}
	return objects
}
*/
