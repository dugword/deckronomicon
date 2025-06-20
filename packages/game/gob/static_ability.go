package gob

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"encoding/json"
	"strings"
)

type StaticAbility struct {
	name      mtg.StaticKeyword
	cost      cost.Cost
	modifiers json.RawMessage `json:"Modifiers,omitempty"`
}

func NewStaticAbility(name mtg.StaticKeyword, modifiers json.RawMessage) StaticAbility {
	staticAbility := StaticAbility{
		name:      name,
		modifiers: modifiers,
	}
	return staticAbility
}

func (a StaticAbility) Name() string {
	return string(a.name)
}

func (a StaticAbility) StaticKeyword() mtg.StaticKeyword {
	return a.name
}

// Description returns a string representation of the static ability.
func (a StaticAbility) Description() string {
	var descriptions []string
	return strings.Join(descriptions, ", ")
}

type SpliceModifiers struct {
	Subtype mtg.Subtype `json:"Subtype,omitempty"`
}
