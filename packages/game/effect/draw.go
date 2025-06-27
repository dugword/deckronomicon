package effect

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Draw struct {
	Count  int
	Target mtg.TargetType
}

func (e Draw) Name() string {
	return "Draw"
}

func NewDraw(modifiers map[string]any) (Draw, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return Draw{}, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return Draw{}, err
	}
	return Draw{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e Draw) TargetSpec() TargetSpec {
	switch e.Target {
	case "", mtg.TargetTypeNone:
		return NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DrawEffect", e.Target))
		return NoneTargetSpec{}
	}
}
