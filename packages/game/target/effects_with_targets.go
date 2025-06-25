package target

import (
	"deckronomicon/packages/game/definition"
	"fmt"
)

type EffectWithTarget struct {
	EffectSpec definition.EffectSpec
	Target     TargetValue
	SourceID   string
}

type EffectTargetKey struct {
	SourceID    string
	EffectIndex int
}

func BuildEffectWithTargets(
	sourceID string,
	effectSpecs []definition.EffectSpec,
	targetsForEffects map[EffectTargetKey]TargetValue,
) ([]EffectWithTarget, error) {
	var effectWithTargets []EffectWithTarget
	for i, effectSpec := range effectSpecs {
		targetsForEffect, ok := targetsForEffects[EffectTargetKey{
			SourceID:    sourceID,
			EffectIndex: i,
		}]
		if !ok {
			return nil, fmt.Errorf("missing targets for effect %d of card %q", i, sourceID)
		}
		effectWithTargets = append(effectWithTargets, EffectWithTarget{
			EffectSpec: effectSpec,
			Target:     targetsForEffect,
			SourceID:   sourceID,
		})
	}
	return effectWithTargets, nil
}
