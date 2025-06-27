package effect

import (
	"fmt"
)

type EffectWithTarget struct {
	Effect   Effect
	Target   Target
	SourceID string
}

type EffectTargetKey struct {
	SourceID    string
	EffectIndex int
}

func BuildEffectWithTargets(
	sourceID string,
	effects []Effect,
	targetsForEffects map[EffectTargetKey]Target,
) ([]EffectWithTarget, error) {
	var effectWithTargets []EffectWithTarget
	for i, effect := range effects {
		targetsForEffect, ok := targetsForEffects[EffectTargetKey{
			SourceID:    sourceID,
			EffectIndex: i,
		}]
		if !ok {
			return nil, fmt.Errorf("missing targets for effect %d of card %q", i, sourceID)
		}
		effectWithTargets = append(effectWithTargets, EffectWithTarget{
			Effect:   effect,
			Target:   targetsForEffect,
			SourceID: sourceID,
		})
	}
	return effectWithTargets, nil
}
