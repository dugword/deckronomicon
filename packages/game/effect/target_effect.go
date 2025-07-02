package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
)

type TargetEffect struct {
	Target mtg.TargetType
}

func NewTarget(modifiers map[string]any) (*TargetEffect, error) {
	targetPermanentModifier, err := parseTargetPermanent(modifiers)
	if err != nil {
		return nil, err
	}
	return &TargetEffect{
		Target: targetPermanentModifier,
	}, nil
}

func (t *TargetEffect) Name() string {
	return "TargetEffect"
}

func (t *TargetEffect) TargetSpec() target.TargetSpec {
	return target.PermanentTargetSpec{}
}
