package effect

import (
	"deckronomicon/packages/game/mtg"
)

type Search struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func NewSearch(modifiers map[string]any) (Search, error) {
	query, err := parseQuery(modifiers)
	if err != nil {
		return Search{}, err
	}
	return Search(query), nil
}

func (e Search) Name() string {
	return "Search"
}

func (e Search) TargetSpec() TargetSpec {
	return NoneTargetSpec{}
}
