package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

// Target is mostly used for unit tests and debugging puroses when a simple
// effect is needed that targets a specific object.
// However, it is an actual effect on "Indicate" :)
type TargetEffect struct {
	Target string
}

func NewTargetEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var targetEffect TargetEffect
	targetStr, ok := effectSpec.Modifiers["Target"].(string)
	if !ok {
		return nil, fmt.Errorf("TargetEffect requires a 'Target' modifier of type string, got %T", effectSpec.Modifiers["Target"])
	}
	targetEffect.Target = targetStr
	return targetEffect, nil
}

func (d TargetEffect) Name() string {
	return "Target Permanent or Player"
}

func (d TargetEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		return target.NoneTargetSpec{}
	case "Permanent":
		return target.PermanentTargetSpec{}
	case "NonlandPermanent":
		panic("NonlandPermanent target spec is not yet implemented for TargetEffect")
		return target.NoneTargetSpec{}
	case "Player":
		return target.PlayerTargetSpec{}
	case "Spell":
		return target.SpellTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for TargetEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e TargetEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	var events []event.GameEvent
	return EffectResult{
		Events: events,
	}, nil
}
