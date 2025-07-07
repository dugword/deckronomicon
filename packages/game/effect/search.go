package effect

import (
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
)

type Search struct {
	QueryOpts query.Opts
}

func NewSearch(modifiers map[string]any) (*Search, error) {
	queryOpts, err := buildQueryOpts(modifiers)
	if err != nil {
		return nil, err
	}
	search := Search{
		QueryOpts: queryOpts,
	}
	return &search, nil
}

func (e *Search) Name() string {
	return "Search"
}

func (e *Search) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
