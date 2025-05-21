package game

import (
	"fmt"
	"strings"
)

// TODO: Implement mana abilities to not use the stack

// Ability is the general interface for all abilities.
type Ability interface {
	Description() string
	Resolve(game *GameState, resolver ChoiceResolver) error
}

// AbilityTag represents tags associated with abilities.
type AbilityTag struct {
	Key   AbilityTagKey
	Value string
}

// ActivatedAbility represents abilities that require activation costs.
type ActivatedAbility struct {
	Cost    Cost
	Effects []*Effect
}

// Description returns a string representation of the activated ability.
func (a *ActivatedAbility) Description() string {
	var effectDescriptions []string
	for _, effect := range a.Effects {
		effectDescriptions = append(effectDescriptions, effect.Description)
	}
	return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(effectDescriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
func (a *ActivatedAbility) IsManaAbility() bool {
	for _, tag := range a.Tags() {
		if tag.Key == AbilityTagManaSource {
			return true
		}
	}
	return false
}

// Resolve resolves the activated ability by paying its cost and applying its
// effects.
func (a *ActivatedAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return fmt.Errorf("cannot resolve effect: err")
		}
	}
	return nil
}

// Tags returns the tags associated with the activated ability.
func (a *ActivatedAbility) Tags() []AbilityTag {
	var tags []AbilityTag
	for _, effect := range a.Effects {
		tags = append(tags, effect.Tags...)
	}
	return tags
}

// SpellAbility represents abilities on instant or sorcery spells.
type SpellAbility struct {
	// Cost    Cost // TODO: Additional costs?
	Effects []*Effect
}

// Description returns a string representation of the spell ability.
func (a *SpellAbility) Description() string {
	var effectDescriptions []string
	for _, effect := range a.Effects {
		effectDescriptions = append(effectDescriptions, effect.Description)
	}
	// Additional Costs
	// return fmt.Sprintf("%s: %s", a.Cost.ToString(), strings.Join(effectDescriptions, ", "))
	return strings.Join(effectDescriptions, ", ")
}

// Resolve resolves the spell ability by applying its effects.
func (a *SpellAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return err
		}
	}
	return nil
}

// StaticAbility represents continuous effects.
type StaticAbility struct {
	Effects []Effect
}

// Description returns a string representation of the static ability.
func (a *StaticAbility) Description() string {
	var effectDescriptions []string
	for _, effect := range a.Effects {
		effectDescriptions = append(effectDescriptions, effect.Description)
	}
	return strings.Join(effectDescriptions, ", ")
}

// TriggeredAbility represents abilities that trigger on specific events.
type TriggeredAbility struct {
	// Cost Cost // TODO: Additoonal Cost
	Effects          []Effect
	TriggerCondition func(event Event) bool
}

// Description returns a string representation of the triggered ability.
func (a *TriggeredAbility) Description() string {
	var effectDescriptions []string
	for _, effect := range a.Effects {
		effectDescriptions = append(effectDescriptions, effect.Description)
	}
	return strings.Join(effectDescriptions, ", ")
}

// Resolve resolves the triggered ability by applying its effects.
func (a *TriggeredAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return fmt.Errorf("cannot resolve effect: err")
		}
	}
	return nil
}
