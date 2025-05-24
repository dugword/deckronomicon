package game

import (
	"fmt"
	"strings"
)

// TODO: Figure this out
// Abilities are Game Objects too....

// TODO: Is this useful?
// Ability is the general interface for all abilities.
/*
type Ability interface {
	Description() string
	Resolve(game *GameState, resolver ChoiceResolver) error
}
*/

const (
	SpeedInstant string = "Instant"
	SpeedSorcery string = "Instant"
)

const (
	AbilityTagDiscard     string = "Discard"
	AbilityTagDraw        string = "Draw"
	AbilityTagManaAbility string = "ManaAbility"
	AbilityTagScry        string = "Scry"
)

// ActivatedAbility represents abilities that require activation costs.
type ActivatedAbility struct {
	Cost    Cost
	Effects []*Effect
	ID      string
	Zone    string
	source  GameObject
	Speed   string
}

func (a *ActivatedAbility) CanPlay(state *GameState) bool {
	if a.Speed == SpeedSorcery && !state.CanCastSorcery() {
		return false
	}
	return true
}

// BuildActivatedAbility builds an activated ability from the given
// specification.
func BuildActivatedAbility(spec ActivatedAbilitySpec, source GameObject) (*ActivatedAbility, error) {
	// ZoneBattlefield is the default zone for activated abilities.
	zone := ZoneBattlefield
	if spec.Zone != "" {
		zone = spec.Zone
	}
	ability := ActivatedAbility{
		ID:     GetNextID(),
		source: source,
		Zone:   zone,
	}
	cost, err := NewCost(spec.Cost, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create cost: %w", err)
	}
	ability.Cost = cost
	for _, effectSpec := range spec.EffectSpecs {
		effect, err := BuildEffect(source, effectSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to create effect: %w", err)
		}
		ability.Effects = append(ability.Effects, effect)
	}
	return &ability, nil
}

// Description returns a string representation of the activated ability.
func (a *ActivatedAbility) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	return fmt.Sprintf("%s: %s", a.Cost.Description(), strings.Join(descriptions, ", "))
}

// IsManaAbility checks if the activated ability is a mana ability.
func (a *ActivatedAbility) IsManaAbility() bool {
	for _, tag := range a.Tags() {
		if tag.Key == AbilityTagManaAbility {
			return true
		}
	}
	return false
}

// Resolve resolves the activated ability. Any costs must be paid before
// resolving the ability.
func (a *ActivatedAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}

// Tags returns the tags associated with the activated ability.
func (a *ActivatedAbility) Tags() []EffectTag {
	var tags []EffectTag
	for _, effect := range a.Effects {
		tags = append(tags, effect.Tags...)
	}
	return tags
}

// StaticAbility represents continuous effects.
type StaticAbility struct {
	ID string
	// TODO: I don't like this, it feels close but not quite right.
	// I think tags should mostly be informational and for finding stuff, not
	// impacting the actual effect logic.
	// I also want to parse the JSON config closer to when the card is created
	// to verify things, and not store logic in strings. Maybe having a
	// defined "Cost" field in this struct would be better, but I don't know
	// the full set of possible static abilties.
	// Maybe I need to create a new Modifier tag...  instead of resuing effect
	// tags just because they are similiar. :shrug:
	Modifiers []EffectTag
}

// BuildStaticAbility builds a static ability from the given specification.
func BuildStaticAbility(spec StaticAbilitySpec, source GameObject) (*StaticAbility, error) {
	ability := StaticAbility{
		ID: spec.ID,
	}
	for _, modifer := range spec.Modifiers {
		efectTag := EffectTag{
			Key:   modifer.Key,
			Value: modifer.Value,
		}
		ability.Modifiers = append(ability.Modifiers, efectTag)
	}
	return &ability, nil
}

// Description returns a string representation of the static ability.
func (a *StaticAbility) Description() string {
	var descriptions []string
	for _, modifier := range a.Modifiers {
		descriptions = append(
			descriptions,
			fmt.Sprintf("%s: %s", modifier.Key, modifier.Value),
		)
	}
	return strings.Join(descriptions, ", ")
}

// SpellAbility represents abilities on instant or sorcery spells.
type SpellAbility struct {
	// Cost    Cost // TODO: Additional costs?
	Effects []*Effect
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
func (a *SpellAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}

// TriggeredAbility represents abilities that trigger on specific events.
type TriggeredAbility struct {
	// Cost Cost // TODO: Additional Cost
	Effects          []Effect
	TriggerCondition func(event Event) bool
}

// Description returns a string representation of the triggered ability.
func (a *TriggeredAbility) Description() string {
	var descriptions []string
	for _, effect := range a.Effects {
		descriptions = append(descriptions, effect.Description)
	}
	return strings.Join(descriptions, ", ")
}

// Resolve resolves the triggered ability by applying its effects.
func (a *TriggeredAbility) Resolve(state *GameState, resolver ChoiceResolver) error {
	for _, effect := range a.Effects {
		if err := effect.Apply(state, resolver); err != nil {
			return fmt.Errorf("cannot resolve effect: %w", err)
		}
	}
	return nil
}
