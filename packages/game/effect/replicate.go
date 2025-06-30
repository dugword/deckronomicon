package effect

import "deckronomicon/packages/game/target"

type Replicate struct {
	Count int `json:"Count,omitempty"`
}

func NewReplicate(modifiers map[string]any) (Replicate, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return Replicate{}, err
	}
	return Replicate{
		Count: countModifier,
	}, nil
}

func (e Replicate) Name() string {
	return "Replicate"
}

func (e Replicate) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
