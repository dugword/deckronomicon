package effect

import (
	"deckronomicon/packages/game/definition"
	"fmt"
)

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
}
