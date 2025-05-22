package game

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

// AddPermanent adds a permanent to the battlefield.
func (b *Battlefield) AddPermanent(permanent *Permanent) {
	b.permanents = append(b.permanents, permanent)
}

// Permanents returns all permanents on the battlefield.
// TODO: For display remove later
// I don't remember why I said for display only... maybe permanents should
// only be manipulated through the methods of this struct.
// Maybe return choices?
func (b *Battlefield) Permanents() []*Permanent {
	return b.permanents
}

// GetPermanents returns the permanent at the given index.
// TODO: need to figure out how to handle duplicates so this can be better,
// maybe by an ID?
func (b *Battlefield) GetPermanent(i int) *Permanent {
	return b.permanents[i]
}

// GetPermanentsWithActivatedAbilities returns a list of permanents with
// activated abilities that can be paid.
func (b *Battlefield) GetPermanentsWithActivatedAbilities(state *GameState) []Choice {
	var found []Choice
	for i, permanent := range b.permanents {
		for _, activatedAbility := range permanent.ActivatedAbilities() {
			if activatedAbility.Cost.CanPay(state) {
				found = append(
					found,
					Choice{
						Name:  permanent.Name(),
						Index: i,
						Zone:  "Battlefield",
					},
				)
				break
			}
		}
	}
	return found
}

// GetUnTappedPermanents returns a list of choices for untapped permanents on
// the battlefield.
func (b *Battlefield) GetUnTappedPermanents(state *GameState) []Choice {
	var found []Choice
	for i, permanent := range b.permanents {
		if !permanent.IsTapped() {
			found = append(found, Choice{Name: permanent.Name(), Index: i})
		}
	}
	return found
}

// GetTappedPermanents returns a list of choices for tapped permanents on
// the battlefield.
func (b *Battlefield) GetTappedPermanents(state *GameState) []Choice {
	var found []Choice
	for i, permanent := range b.permanents {
		if permanent.IsTapped() {
			found = append(found, Choice{Name: permanent.Name(), Index: i})
		}
	}
	return found
}

// FindUntappedPermanent returns the first untapped permanent with the given
// name.
// TODO: this only finds the first one, that will be problematic if there are
// duplicates where it matters
func (b *Battlefield) FindUntappedPermanent(name string) *Permanent {
	for _, p := range b.permanents {
		if p.Name() == name && !p.IsTapped() {
			return p
		}
	}
	return nil
}

// FindTappedPermanent returns the first tapped permanent with the given name.
// TODO: this only finds the first one, that will be problematic if there are
// duplicates where it matters
func (b *Battlefield) FindTappedPermanent(name string) *Permanent {
	for _, p := range b.permanents {
		if p.Name() == name && p.IsTapped() {
			return p
		}
	}
	return nil
}

// FindPermanentWithAvailableActivatedAbility returns the first permanent with
// an activated ability that can be paid.
func (b *Battlefield) FindPermanentWithAvailableActivatedAbility(name string, state *GameState) *Permanent {
	for _, permanent := range b.permanents {
		for _, activatedAbility := range permanent.ActivatedAbilities() {
			if activatedAbility.Cost.CanPay(state) {
				return permanent
			}
		}
	}
	return nil
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

// RemoveSummoningSickness removes summoning sickness from all permanents
func (b *Battlefield) RemoveSummoningSickness() {
	for _, permanent := range b.permanents {
		permanent.RemoveSummoningSickness()
	}
}
