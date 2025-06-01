package activated

import (
	"deckronomicon/packages/game/core"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
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

type Source interface {
	Name() string
}

// Ability represents abilities that require activation costs.
type Ability struct {
	name    string
	Cost    cost.Cost
	Effects []*effect.Effect
	id      string
	Zone    string
	source  Source
	Speed   string
}

// BuildActivatedAbility builds an activated ability from the given
// specification.
// TODO use a better interface
func BuildActivatedAbility(state core.State, spec definition.ActivatedAbilitySpec, source query.Object) (*Ability, error) {
	// ZoneBattlefield is the default zone for activated abilities.
	/*
		zone := ZoneBattlefield
		if spec.Zone != "" {
			zone = spec.Zone
		}
	*/
	ability := Ability{
		name:   spec.Name,
		id:     state.GetNextID(),
		source: source,
		Zone:   "",
		Speed:  spec.Speed,
	}
	cost, err := cost.NewCost(spec.Cost, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create cost: %w", err)
	}
	ability.Cost = cost
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := effect.BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}

func (a *Ability) Name() string {
	return fmt.Sprintf("%s - %s", a.source.Name(), a.name)
}

func (a *Ability) ID() string {
	return a.id
}

func (a *Ability) CanActivate(state core.State, playerID string) bool {
	if a.Speed == SpeedSorcery && !state.CanCastSorcery(playerID) {
		return false
	}
	return true
}

// CardTypese returns the card types associated with the activated ability.
// This exists to satisfy the Object interface, but activated abilities
// typically do not have card types.
func (a *Ability) CardTypes() []mtg.CardType {
	return nil
}

// Subtypese returns the card types associated with the activated ability.
// This exists to satisfy the Object interface, but activated abilities
// typically do not have card types.
func (a *Ability) Subtypes() []mtg.Subtype {
	return nil
}

// Supertypes returns the card types associated with the activated ability.
// This exists to satisfy the Object interface, but activated abilities
// typically do not have card types.
func (a *Ability) Supertypes() []mtg.Supertype {
	return nil
}

func (a *Ability) Colors() mtg.Colors {
	// Activated abilities typically do not have colors, but this method
	// exists to satisfy the Card interface.
	return mtg.Colors{}
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
		if tag.Key == effect.TagManaAbility {
			return true
		}
	}
	return false
}

func (a *Ability) Match(p query.Predicate) bool {
	return p(a)
}

// Resolve resolves the activated ability. Any costs must be paid before
// resolving the ability.
func (a *Ability) Resolve(state core.State, player core.Player) error {
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
