package effect

import "deckronomicon/packages/game/mtg"

type TargetEffect struct {
	Target mtg.TargetType
}

func NewTarget(modifiers map[string]any) (TargetEffect, error) {
	targetPermanentModifier, err := parseTargetPermanent(modifiers)
	if err != nil {
		return TargetEffect{}, err
	}
	return TargetEffect{
		Target: targetPermanentModifier,
	}, nil
}

func (t TargetEffect) Name() string {
	return "TargetEffect"
}

func (t TargetEffect) TargetSpec() TargetSpec {
	return PermanentTargetSpec{}
}
