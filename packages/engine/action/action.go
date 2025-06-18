package action

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/target"
	"fmt"
)

type EffectTargetKey struct {
	SourceID    string
	EffectIndex int
}

func buildEffectWithTargets(
	sourceID string,
	effectSpecs []definition.EffectSpec,
	targetsForEffects map[EffectTargetKey]target.TargetValue,
) ([]gob.EffectWithTarget, error) {
	var effectWithTargets []gob.EffectWithTarget
	for i, effectSpec := range effectSpecs {
		targetsForEffect, ok := targetsForEffects[EffectTargetKey{
			SourceID:    sourceID,
			EffectIndex: i,
		}]
		if !ok {
			return nil, fmt.Errorf("missing targets for effect %d of card %q", i, sourceID)
		}
		effectWithTargets = append(effectWithTargets, gob.EffectWithTarget{
			EffectSpec: effectSpec,
			Target:     targetsForEffect,
			SourceID:   sourceID,
		})
	}
	return effectWithTargets, nil
}
