package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type Discard struct {
	Count  int
	Target mtg.TargetType
}

func (e *Discard) Name() string {
	return "Discard"
}

func NewDiscard(modifiers map[string]any) (*Discard, error) {
	countModifier, err := parseCount(modifiers)
	if err != nil {
		return nil, err
	}
	targetPlayerModifier, err := parseTargetPlayer(modifiers)
	if err != nil {
		return nil, err
	}
	return &Discard{
		Count:  countModifier,
		Target: targetPlayerModifier,
	}, nil
}

func (e *Discard) TargetSpec() target.TargetSpec {
	switch e.Target {
	case "", mtg.TargetTypeNone:
		return target.NoneTargetSpec{}
	case mtg.TargetTypePlayer:
		return target.PlayerTargetSpec{}
	default:
		panic(fmt.Sprintf("unknown target spec %q for DiscardEffect", e.Target))
		return target.NoneTargetSpec{}
	}
}
