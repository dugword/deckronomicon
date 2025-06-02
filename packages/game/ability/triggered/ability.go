package triggered

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/effect"
	"fmt"
	"strings"
)

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
func (a *Ability) Resolve(state core.State, player core.Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
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

/*
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
*/
