package effect

import (
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type AdditionalManaEffect struct {
	Subtype  string `json:"Subtype"`
	Mana     string `json:"Mana"`
	Duration string `json:"Duration"`
}

func NewAdditionalManaEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var additionalManaEffect AdditionalManaEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &additionalManaEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AdditionalManaEffect: %w", err)
	}
	return additionalManaEffect, nil
}

func (e AdditionalManaEffect) Name() string {
	return "AdditionalMana"
}

func (e AdditionalManaEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e AdditionalManaEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {

	return EffectResult{
		Events: nil,
	}, nil
}

func (e AdditionalManaEffect) ResolveNew(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
	resEnv *resenv.ResEnv,
) (EffectResult, error) {

	return EffectResult{
		Events: nil,
	}, nil
}
