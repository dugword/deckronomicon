package game

import (
	"errors"
	"fmt"
)

// Permanents returns all permanents on the battlefield.
// TODO: For display remove later, permanents should only be manipulated
// through the Battlefield interface.
func (b *Battlefield) Permanents() []*Permanent {
	return b.permanents
}

// Battlefield represents the battlefield during an active game.
type Battlefield struct {
	permanents []*Permanent
}

// NewBattlefield creates a new Battlefield instance.
func NewBattlefield() *Battlefield {
	battlefield := Battlefield{
		permanents: []*Permanent{},
	}
	return &battlefield
}

func (b *Battlefield) Add(object GameObject) error {
	permanent, ok := object.(*Permanent)
	if !ok {
		return errors.New("object is not a permanent")
	}
	b.permanents = append(b.permanents, permanent)
	return nil
}

func (b *Battlefield) AvailableActivatedAbilities(state *GameState, player *Player) []GameObject {
	var objects []GameObject
	for _, permanent := range b.permanents {
		for _, ability := range permanent.ActivatedAbilities() {
			if !ability.CanPlay(state) {
				continue
			}
			if ability.Cost.CanPay(state, player) {
				objects = append(objects, ability)
			}
		}
	}
	return objects
}

// AvailableToPlay returns a list of permanents that can be played from the
// battlefield. This exists to satisfy the Zone interface. Permanents on the
// battlefield generally cannot be played. There is no rule that prevents a
// card from enabling this, but I don't think any exist.
func (b *Battlefield) AvailableToPlay(*GameState, *Player) []GameObject {
	return nil
}

func (b *Battlefield) Find(id string) (GameObject, error) {
	for _, permanent := range b.permanents {
		if permanent.ID() == id {
			return permanent, nil
		}
	}
	return nil, fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) Get(id string) (GameObject, error) {
	for _, permanent := range b.permanents {
		if permanent.ID() == id {
			return permanent, nil
		}
	}
	return nil, fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) GetAll() []GameObject {
	var all []GameObject
	for _, permanent := range b.permanents {
		all = append(all, permanent)
	}
	return all
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

func (b *Battlefield) Take(id string) (GameObject, error) {
	for i, permanent := range b.permanents {
		if permanent.ID() == id {
			b.permanents = append(b.permanents[:i], b.permanents[i+1:]...)
			return permanent, nil
		}
	}
	return nil, fmt.Errorf("permanent with ID %s not found", id)
}

func (b *Battlefield) Size() int {
	return len(b.permanents)
}

func (b *Battlefield) ZoneType() string {
	return ZoneBattlefield
}

// Battlefield Specific Methods

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
