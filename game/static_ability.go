package game

import (
	"fmt"
	"strings"
)

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
