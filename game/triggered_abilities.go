package game

import (
	"fmt"
	"strconv"
	"strings"
)

// TriggeredAbility represents abilities that trigger on specific events.
type TriggeredAbility struct {
	name string
	id   string
	// Cost Cost // TODO: Additional Cost
	Effects          []Effect
	TriggerCondition func(event Event) bool
}

// TODO
// NewTriggeredAbility

// Description returns a string representation of the triggered ability.
func (a *TriggeredAbility) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	return strings.Join(descriptions, ", ")
}

// Resolve resolves the triggered ability by applying its effects.
func (a *TriggeredAbility) Resolve(state *GameState, player *Player) error {
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

func (a *TriggeredAbility) HasSubtype(Subtype) bool {
	return false
}
func (a *TriggeredAbility) Name() string {
	return a.name
}

func (a *TriggeredAbility) ActivatedAbilities() []*ActivatedAbility {
	return nil
}

func (a *TriggeredAbility) StaticAbilities() []*StaticAbility {
	return nil
}

func (a *TriggeredAbility) HasStaticAbility(string) bool {
	return false
}

func (a *TriggeredAbility) ID() string {
	return a.id
}

func BuildReplicateAbility(card *Card, replicateCount int) *TriggeredAbility {
	handler := func(state *GameState, player *Player) error {
		for range replicateCount {
			spell, err := NewSpell(card)
			if err != nil {
				return fmt.Errorf("failed to create spell from %s: %w", card.Name(), err)
			}
			state.Stack.Add(spell)
		}
		return nil
	}
	tags := []EffectTag{{Key: "Count", Value: strconv.Itoa(replicateCount)}}
	description := fmt.Sprintf("Replicate %s %d times", card.Name(), replicateCount)
	replicateAbility := TriggeredAbility{
		name: "Replicate Trigger",
		id:   GetNextID(),
		Effects: []Effect{
			{
				Apply:       handler,
				Tags:        tags,
				Description: description,
			},
		},
	}
	return &replicateAbility
}
