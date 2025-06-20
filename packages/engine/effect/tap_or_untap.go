package effect

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type TapOrUntapEffect struct {
	Target string `json:"Target"`
}

func NewTapOrUntapEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var tapOrUntapEffect TapOrUntapEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &tapOrUntapEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TapOrUntapEffectModifiers: %w", err)
	}
	return tapOrUntapEffect, nil
}

func (d TapOrUntapEffect) Name() string {
	return "Tap or Untap"
}

func (d TapOrUntapEffect) TargetSpec() target.TargetSpec {
	switch d.Target {
	case "":
		return target.NoneTargetSpec{}
	case "Permanent":
		return target.PermanentTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for TapOrUntapEffect", d.Target))
		return target.NoneTargetSpec{}
	}
}

func (e TapOrUntapEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {
	return EffectResult{
		Events: []event.GameEvent{
			event.UntapPermanentEvent{
				PlayerID:    player.ID(),
				PermanentID: target.ObjectID,
			},
		},
	}, nil
}
