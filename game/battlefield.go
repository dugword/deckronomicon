package game

type Battlefield struct {
	permanents []*Permanent
}

func (b *Battlefield) AddPermanent(permanent *Permanent) {
	b.permanents = append(b.permanents, permanent)
}

// TODO: For display remove later
func (b *Battlefield) Permanents() []*Permanent {
	return b.permanents
}

// TODO: need to figure out how to handle duplicates so this can be better, maybe by an ID?
func (b *Battlefield) GetPermanent(i int) *Permanent {
	return b.permanents[i]
}

func (b *Battlefield) GetPermanentsWithActivatedAbilities(state *GameState) []Choice {
	var found []Choice
	for i, permanent := range b.permanents {
		// TODO: Only supports tap costs currently
		for _, activatedAbility := range permanent.ActivatedAbilities() {
			if activatedAbility.Cost.CanPay(state, permanent) {
				found = append(found, Choice{Name: permanent.Name(), Index: i})
				break
			}
		}
	}
	return found
}

// TODO: this only finds the first one, that will be problematic if there are duplicates where it matters
func (b *Battlefield) FindUntappedPermanent(name string) *Permanent {
	for _, p := range b.permanents {
		if p.Name() == name && !p.IsTapped() {
			return p
		}
	}
	return nil
}

func (b *Battlefield) FindPermanentWithAvailableActivatedAbility(name string, state *GameState) *Permanent {
	for _, permanent := range b.permanents {
		// TODO: Support permanenents with more than 1 ability
		for _, activatedAbility := range permanent.card.ActivatedAbilities() {
			if activatedAbility.Cost.CanPay(state, permanent) {
				return permanent
			}
		}
	}
	return nil
}

func (b *Battlefield) UntapPermanents() {
	for _, permanent := range b.permanents {
		permanent.Untap()
	}
}

func (b *Battlefield) RemoveSummoningSickness() {
	for _, permanent := range b.permanents {
		permanent.RemoveSummoningSickness()
	}
}
