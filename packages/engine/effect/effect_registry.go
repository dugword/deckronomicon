package effect

import (
	"deckronomicon/packages/game/definition"
	"fmt"
)

// TODO: Think through how to make this work more like actions? I think maybe
// I need to split the effect into two parts: a "spec" that describes the effect
// and a "runtime" that is the actual effect that can be applied to the game.

type EffectConstructor func(effectSpec definition.EffectSpec) (Effect, error)

var registry = map[string]EffectConstructor{}

func Register(name string, factory EffectConstructor) {
	registry[name] = factory
}

func Build(effectSpec definition.EffectSpec) (Effect, error) {
	factory, ok := registry[effectSpec.Name]
	if !ok {
		return nil, fmt.Errorf("unknown effect: %s", effectSpec.Name)
	}
	effect, err := factory(effectSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to build effect %s: %w", effectSpec.Name, err)
	}
	return effect, nil
}

// TODO: I don't like this
func init() {
	Register("AddMana", NewAddManaEffect)
	Register("Draw", NewDrawEffect)
	Register("Counterspell", NewCounterspellEffect)
	Register("Scry", NewScryEffect)
	Register("PutBackOnTop", NewPutBackOnTopEffect)
	Register("Search", NewSearchEffect)
	Register("TapOrUntap", NewTapOrUntapEffect)
	Register("AdditionalMana", NewAdditionalManaEffect)
	Register("Tap", NewTapEffect)
	Register("Discard", NewDiscardEffect)
	Register("LookAndChoose", NewLookAndChooseEffect)
	Register("Mill", NewMillEffect)
	Register("ShuffleFromGraveyard", NewShuffleFromGraveyardEffect)
	Register("RegisterDelayedEffect", NewRegisterDelayedEffectEffect)
	Register("Replicate", NewReplicateEffect)
	Register("Target", NewTargetEffect)
}
