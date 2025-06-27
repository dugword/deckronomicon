package effect

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Mill struct {
	Count  int
	Target mtg.TargetType
}

func NewMill(modifiers map[string]any) (Mill, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return Mill{}, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return Mill{}, err
	}
	return Mill{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e Mill) Name() string {
	return "Mill"
}

func (e Mill) TargetSpec() TargetSpec {
	switch e.Target {
	case mtg.TargetTypeNone:
		return NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for MillEffect", e.Target))
		return NoneTargetSpec{}
	}
}
