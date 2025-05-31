package activated

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"fmt"
	"strings"
)

// These are shared constants, would be good to make them type safe and in a
// common package.
const (
	// SpeedSorcery represents the speed of a sorcery speed ability.
	SpeedSorcery = "Sorcery"
	// SpeedInstant represents the speed of an instant speed ability.
	SpeedInstant = "Instant"
)

// TODO: Is this an effect tag?
const AbilityTagManaAbility = "ManaAbility"

// TODO: Add something here so only GameState can be pass in.
type State any

// TODO: Add something here so only Player can be passed in.
type Player any

type Source interface {
	Name() string
}

type SpeedChecker interface {
	CanCastSorcery() bool
}

// Ability represents abilities that require activation costs.
type Ability struct {
	name    string
	Cost    cost.Cost
	Effects []effect.Effect
	id      string
	Zone    string
	Source  Source
	Speed   string
}

func (a *Ability) Name() string {
	return fmt.Sprintf("%s - %s", a.Source.Name(), a.name)
}

func (a *Ability) ID() string {
	return a.id
}

func (a *Ability) CanActivate(checker SpeedChecker) bool {
	if a.Speed == SpeedSorcery && !checker.CanCastSorcery() {
		return false
	}
	return true
}

// Description returns a string representation of the activated ability.
func (a *Ability) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description())
	}
	return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(descriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
func (a *Ability) IsManaAbility() bool {
	for _, tag := range a.Tags() {
		if tag.Key == AbilityTagManaAbility {
			return true
		}
	}
	return false
}

// Resolve resolves the activated ability. Any costs must be paid before
// resolving the ability.
func (a *Ability) Resolve(state State, player Player) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, player); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}

// Tags returns the tags associated with the activated ability.
func (a *Ability) Tags() []effect.Tag {
	var tags []effect.Tag
	for _, effect := range a.Effects {
		tags = append(tags, effect.Tags()...)
	}
	return tags
}
