package effect

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Discard struct {
	Count  int
	Target mtg.TargetType
}

func (e Discard) Name() string {
	return "Discard"
}

func NewDiscard(modifiers map[string]any) (Discard, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return Discard{}, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return Discard{}, err
	}
	return Discard{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e Discard) TargetSpec() TargetSpec {
	switch e.Target {
	case "", mtg.TargetTypeNone:
		return NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DiscardEffect", e.Target))
		return NoneTargetSpec{}
	}
}
