package effect

import (
	"deckronomicon/packages/game/mtg"
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

func (e Counterspell) TargetSpec() TargetSpec {
	return SpellTargetSpec(
		e,
	)
}
