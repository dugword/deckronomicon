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

type TapEffect struct {
	Target string `json:"Target"`
}

func NewTapEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var tapEffect TapEffect
	targetStr, ok := effectSpec.Modifiers["Target"].(string)
	if !ok {
		return nil, fmt.Errorf("TapEffect requires a 'Target' modifier of type string, got %T", effectSpec.Modifiers["Target"])
	}
	if targetStr != "" && targetStr != "Permanent" {
		return nil, fmt.Errorf("TapEffect requires a 'Target' modifier of either empty or 'Permanent', got %q", targetStr)
	}
	return tapEffect, nil
}

func (d TapEffect) Name() string {
	return "Tap or Untap"
}

func (d TapEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		return target.NoneTargetSpec{}
	case "Permanent":
		return target.PermanentTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for TapEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e TapEffect) Resolve(
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
