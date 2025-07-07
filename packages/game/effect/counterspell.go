package effect

import (
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
)

type Counterspell struct {
	QueryOpts query.Opts
}

func (e *Counterspell) Name() string {
	return "Counterspell"
}

func NewCounterspell(modifiers map[string]any) (*Counterspell, error) {
	queryOpts, err := buildQueryOpts(modifiers)
	if err != nil {
		return nil, err
	}
	counterspell := Counterspell{
		QueryOpts: queryOpts,
	}
	return &counterspell, nil
}

func (e *Counterspell) TargetSpec() target.TargetSpec {
	return target.SpellTargetSpec{
		QueryOpts: e.QueryOpts,
	}
}
