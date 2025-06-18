package effect

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

type ShuffleFromGraveyardEffect struct {
	Mana string `json:"Mana"`
}

func NewShuffleFromGraveyardEffect(effectSpec definition.EffectSpec) (Effect, error) {
	var shuffleFromGraveyardEffect ShuffleFromGraveyardEffect
	if err := json.Unmarshal(effectSpec.Modifiers, &shuffleFromGraveyardEffect); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ShuffleFromGraveyardEffect: %w", err)
	}
	return shuffleFromGraveyardEffect, nil
}

func (e ShuffleFromGraveyardEffect) Name() string {
	return "ShuffleFromGraveyard"
}

func (e ShuffleFromGraveyardEffect) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}

func (e ShuffleFromGraveyardEffect) Resolve(
	game state.Game,
	player state.Player,
	source query.Object,
	target target.TargetValue,
) (EffectResult, error) {
	return EffectResult{}, nil
}
