package effect

import (
	"deckronomicon/packages/game/target"
)

type GainLife struct {
	Count int
}

func (e *GainLife) Name() string {
	return "GainLife"
}

func NewGainLife(modifiers map[string]any) (*GainLife, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return nil, err
	}
	return &GainLife{
		Count: countModifier,
	}, nil
}

func (e *GainLife) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
