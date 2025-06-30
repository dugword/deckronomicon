package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
)

type Counterspell struct {
	CardTypes  []mtg.CardType `json:"CardTypes,omitempty"`
	Colors     []mtg.Color    `json:"Colors,omitempty"`
	Subtypes   []mtg.Subtype  `json:"Subtypes,omitempty"`
	ManaValues []int          `json:"ManaValues,omitempty"`
}

func (e Counterspell) Name() string {
	return "Counterspell"
}

func NewCounterspell(modifiers map[string]any) (Counterspell, error) {
	query, err := parseQuery(modifiers)
	if err != nil {
		return Counterspell{}, err
	}
	return Counterspell(query), nil
}

func (e Counterspell) TargetSpec() target.TargetSpec {
	return target.SpellTargetSpec(
		e,
	)
}
