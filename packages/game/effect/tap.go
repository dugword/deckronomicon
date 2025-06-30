package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type Tap struct {
	Target mtg.TargetType `json:"Target"`
}

func NewTap(modifiers map[string]any) (Tap, error) {
	targetPermanentModifier, err := parseTargetPermanent(modifiers)
	if err != nil {
		return Tap{}, err
	}
	return Tap{
		Target: targetPermanentModifier,
	}, nil
}

func (t Tap) Name() string {
	return "Tap"
}

func (e Tap) TargetSpec() target.TargetSpec {
	switch e.Target {
	case mtg.TargetTypePermanent:
		return target.PermanentTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for TapEffect", e.Target))
		return target.NoneTargetSpec{}
	}
}
