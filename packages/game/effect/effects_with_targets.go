package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type EffectWithTarget struct {
	Effect   Effect
	Target   target.Target
	SourceID string
}

type EffectTargetKey struct {
	SourceID    string
	EffectIndex int
}

func BuildEffectWithTargets(
	sourceID string,
	effects []Effect,
	targetsForEffects map[EffectTargetKey]target.Target,
) ([]EffectWithTarget, error) {
	var effectWithTargets []EffectWithTarget
	for i, effect := range effects {
		targetsForEffect, ok := targetsForEffects[EffectTargetKey{
			SourceID:    sourceID,
			EffectIndex: i,
		}]
		if !ok {
			// TODO: Hacky fix me
			switch effect.TargetSpec().(type) {
			case nil, target.NoneTargetSpec:
				effectWithTargets = append(effectWithTargets, EffectWithTarget{
					Effect:   effect,
					Target:   target.Target{Type: mtg.TargetTypeNone},
					SourceID: sourceID,
				})
				continue
			}
			// TODO: Figure out how to not have to specify targets when the effect
			// doesn't require them.
			// TODO: Maybe use that interface thing we did with costs with targets
			// Have the effect struct have a method that returns whether it needs targets
			// or not. And use interfaces type checking to narrow the type.
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
