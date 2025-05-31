package spell

import (
	"deckronomicon/packages/game/effect"
	"fmt"
	"strings"
)

type State interface {
	Players() []Player
}

type Player interface {
	Agent() Agent
}

type Agent interface {
	ReportState(State)
}

// Ability represents abilities on instant or sorcery spells.
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
func (a *Ability) Resolve(state State, player Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
		for _, p := range state.Players() {
			p.Agent().ReportState(state)
		}
	}
	return nil
}
