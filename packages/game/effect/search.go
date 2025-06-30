package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
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

func (e Search) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
