package game

import (
	"fmt"
	"strings"
)

// SpellAbility represents abilities on instant or sorcery spells.
type SpellAbility struct {
	// Cost    Cost // TODO: Additional costs?
	Effects []*Effect
}

// Splice
func (s *SpellAbility) Splice(spell *SpellAbility) {
	if spell == nil {
		return
	}
	s.Effects = append(s.Effects, spell.Effects...)
}

// BuildSpellAbility builds a spell ability from the given specification.
func BuildSpellAbility(spec *SpellAbilitySpec, source GameObject) (*SpellAbility, error) {
	ability := SpellAbility{}
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}

// Description returns a string representation of the spell ability.
func (a *SpellAbility) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	return strings.Join(descriptions, ", ")
}

// Resolve resolves the spell ability by applying its effects.
func (a *SpellAbility) Resolve(state *GameState, player *Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
		for _, p := range state.Players {
			p.Agent.ReportState(state)
		}
	}
	return nil
}
