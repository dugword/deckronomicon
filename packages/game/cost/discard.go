package cost

import "deckronomicon/packages/game/target"

type DiscardThis struct{}

func (c DiscardThis) isCost() {}

func (c DiscardThis) Description() string {
	return "Discard this card"
}

type DiscardACard struct {
	target target.Target
}

func (c DiscardACard) isCost() {}

func (c DiscardACard) Description() string {
	return "Discard a card"
}

func (c DiscardACard) Target() target.Target {
	return c.target
}

func (c DiscardACard) WithTarget(target target.Target) CostWithTarget {
	c.target = target
	return c
}

func (c DiscardACard) TargetSpec() target.TargetSpec {
	return target.CardTargetSpec{}
}
