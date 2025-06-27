package effect

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type TapOrUntap struct {
	Target mtg.TargetType `json:"Target"`
}

func NewTapOrUntap(modifiers map[string]any) (TapOrUntap, error) {
	targetPermanentModifier, err := parseTargetPermanent(modifiers)
	if err != nil {
		return TapOrUntap{}, err
	}
	return TapOrUntap{
		Target: targetPermanentModifier,
	}, nil
}

func (e TapOrUntap) Name() string {
	return "TapOrUntap"
}

func (e TapOrUntap) TargetSpec() TargetSpec {
	switch e.Target {
	case mtg.TargetTypePermanent:
		return PermanentTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for TapOrUntapEffect", e.Target))
		return NoneTargetSpec{}
	}
}
