package gob

import (
	// "deckronomicon/packages/game/definition"

	"deckronomicon/packages/game/mtg"
	"encoding/json"
	"fmt"
	"strings"
)

// StaticAbility represents continuous effects.
type StaticAbility struct {
	// TODO: Maybe add a typed "Keyword" value here
	// keyword mtg.keyword
	name mtg.StaticKeyword
	// TODO: I don't like this, it feels close but not quite right.
	// I think tags should mostly be informational and for finding stuff, not
	// impacting the actual effect logic.
	// I also want to parse the JSON config closer to when the card is created
	// to verify things, and not store logic in strings. Maybe having a
	// defined "Cost" field in this struct would be better, but I don't know
	// the full set of possible static abilties.
	// Maybe I need to create a new Modifier tag...  instead of resuing effect
	// tags just because they are similiar. :shrug:
	// TODO: Or maybe this is fine.
	Modifiers json.RawMessage `json:"Modifiers,omitempty"`
}

/*
	func NewStaticAbility(id string, modifiers []Tag) StaticAbility {
		staticAbility := StaticAbility{
			id:        id,
			Modifiers: modifiers,
		}
		return staticAbility
	}
*/
func (a StaticAbility) Name() string {
	return string(a.name)
}

func (a StaticAbility) StaticKeyword() mtg.StaticKeyword {
	keyword, ok := mtg.StringToStaticKeyword(string(a.name))
	if !ok {
		panic(fmt.Sprintf("invalid static keyword: %s", a.name))
	}
	return keyword
}

// Description returns a string representation of the static ability.
func (a StaticAbility) Description() string {
	var descriptions []string
	/*
		for _, modifier := range a.Modifiers {
			descriptions = append(
				descriptions,
				fmt.Sprintf("%s: %s", modifier.Key, modifier.Value),
			)
		}
	*/
	return strings.Join(descriptions, ", ")
}

/*
// BuildStaticAbility builds a static ability from the given specification.
func BuildStaticAbility(spec definition.StaticAbilitySpec, source query.Object) (*StaticAbility, error) {
	ability := StaticAbility{
		// TODO: Use string types
		id: string(spec.ID),
	}
	for _, modifer := range spec.Modifiers {
		efectTag := Tag{
			Key:   modifer.Key,
			Value: modifer.Value,
		}
		ability.Modifiers = append(ability.Modifiers, efectTag)
	}
	return &ability, nil
}
*/
