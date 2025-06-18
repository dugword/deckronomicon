package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type TapEffect struct {
	Target string `json:"Target"`
}

func NewTapEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var tapEffect TapEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &tapEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TapEffectModifiers: %w", err)
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
) (EffectResult, error) {
	var events []event.GameEvent
	return EffectResult{
		Events: events,
	}, nil
}
