package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type Mill struct {
	Count  int
	Target mtg.TargetType
}

func NewMill(modifiers map[string]any) (*Mill, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return nil, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return nil, err
	}
	return &Mill{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e *Mill) Name() string {
	return "Mill"
}

func (e *Mill) TargetSpec() target.TargetSpec {
	switch e.Target {
	case mtg.TargetTypeNone:
		return target.NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for MillEffect", e.Target))
		return target.NoneTargetSpec{}
	}
}
