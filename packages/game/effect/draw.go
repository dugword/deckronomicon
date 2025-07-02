package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type Draw struct {
	Count  int
	Target mtg.TargetType
}

func (e *Draw) Name() string {
	return "Draw"
}

func NewDraw(modifiers map[string]any) (*Draw, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return nil, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return nil, err
	}
	return &Draw{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e *Draw) TargetSpec() target.TargetSpec {
	switch e.Target {
	case "", mtg.TargetTypeNone:
		return target.NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DrawEffect", e.Target))
		return target.NoneTargetSpec{}
	}
}
