package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/add"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/remove"
	"deckronomicon/packages/query/take"
)

// Battlefield represents the battlefield during an active game.
type Battlefield struct {
	permanents []*gob.Permanent
}

// NewBattlefield creates a new Battlefield instance.
func NewBattlefield() *Battlefield {
	battlefield := Battlefield{
		permanents: []*gob.Permanent{},
	}
	return &battlefield
}

func (b *Battlefield) Add(permanents ...*gob.Permanent) *Battlefield {
	return &Battlefield{
		permanents: add.Item(b.permanents, permanents...),
	}
}

func (b *Battlefield) Contains(predicate query.Predicate) bool {
	return query.Contains(b.permanents, predicate)
}

func (b *Battlefield) Find(predicate query.Predicate) (*gob.Permanent, bool) {
	return query.Find(b.permanents, predicate)
}

func (b *Battlefield) FindAll(predicate query.Predicate) []*gob.Permanent {
	return query.FindAll(b.permanents, predicate)
}

func (b *Battlefield) Get(id string) (*gob.Permanent, bool) {
	return query.Get(b.permanents, id)
}

func (b *Battlefield) GetAll() []*gob.Permanent {
	return query.GetAll(b.permanents)
}

func (b *Battlefield) Name() string {
	return string(mtg.ZoneBattlefield)
}

func (b *Battlefield) Remove(id string) (*Battlefield, bool) {
	permanents, ok := remove.By(b.permanents, has.ID(id))
	if !ok {
		return nil, false
	}
	return &Battlefield{permanents: permanents}, true
}

func (b *Battlefield) RemoveBy(predicate query.Predicate) (*Battlefield, bool) {
	permanents, ok := remove.By(b.permanents, predicate)
	if !ok {
		return nil, false
	}
	return &Battlefield{permanents: permanents}, true
}

func (b *Battlefield) Size() int {
	return len(b.permanents)
}

func (b *Battlefield) Take(id string) (*gob.Permanent, *Battlefield, bool) {
	permanent, permanents, ok := take.By(b.permanents, has.ID(id))
	if !ok {
		return nil, nil, false
	}
	return permanent, &Battlefield{permanents: permanents}, true
}

func (b *Battlefield) TakeBy(predicate query.Predicate) (*gob.Permanent, *Battlefield, bool) {
	permanent, permanents, ok := take.By(b.permanents, predicate)
	if !ok {
		return nil, nil, false
	}
	return permanent, &Battlefield{permanents: permanents}, true
}

func (b *Battlefield) UntapAll(playerID string) *Battlefield {
	battlefield := NewBattlefield()
	for _, p := range b.permanents {
		if p.Controller() == playerID {
			p = p.Untap()
		}
		battlefield = battlefield.Add(p)
	}
	return battlefield
}

func (b *Battlefield) RemoveSummoningSickness(playerID string) *Battlefield {
	battlefield := NewBattlefield()
	for _, p := range b.permanents {
		if p.Controller() == playerID {
			p = p.RemoveSummoningSickness()
		}
		battlefield = battlefield.Add(p)
	}
	return battlefield
}

func (b *Battlefield) Zone() mtg.Zone {
	return mtg.ZoneBattlefield
}

func (b *Battlefield) WithUpdatedPermanent(permanent *gob.Permanent) *Battlefield {
	var permanents []*gob.Permanent
	for _, p := range b.permanents {
		if p.ID() == permanent.ID() {
			permanents = append(permanents, permanent)
			continue
		}
		permanents = append(permanents, p)
	}
	return &Battlefield{
		permanents: permanents,
	}
}
