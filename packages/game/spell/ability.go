package spell

// Ability represents abilities on instant or sorcery spells.
/*
type Ability struct {
	// Cost    Cost // TODO: Additional costs?
	Effects []*effect.Effect
}

// Splice
func (s *Ability) Splice(spell *Ability) {
	if spell == nil {
		return
	}
	s.Effects = append(s.Effects, spell.Effects...)
}

// Description returns a string representation of the spell ability.
func (a *Ability) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description())
	}
	return strings.Join(descriptions, ", ")
}

// Resolve resolves the spell ability by applying its effects.
func (a *Ability) Resolve(state core.State, plyr core.Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, plyr); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}
*/
