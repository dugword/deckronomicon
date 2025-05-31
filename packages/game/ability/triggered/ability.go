package triggered

import (
	"deckronomicon/packages/game/effect"
	"fmt"
	"strings"
)

type agent interface {
	ReportState(state)
}

type player interface {
	Agent() agent
}

type state interface {
	GetNextID() string
	Players() []player
}

// Ability represents abilities that trigger on specific events.
type Ability struct {
	name string
	id   string
	// Cost Cost // TODO: Additional Cost
	Effects []effect.Effect
	// TriggerCondition func(event Event) bool
}

// TODO
// NewAbility

// Description returns a string representation of the triggered ability.
func (a *Ability) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description())
	}
	return strings.Join(descriptions, ", ")
}

// Resolve resolves the triggered ability by applying its effects.
func (a *Ability) Resolve(state state, player player) error {
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

func (a *Ability) Name() string {
	return a.name
}

func (a *Ability) ID() string {
	return a.id
}
