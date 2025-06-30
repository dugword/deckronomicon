package cost

import (
	"deckronomicon/packages/game/target"
	"fmt"
)

type SacrificeThis struct{}

func (c SacrificeThis) isCost() {}

func (c SacrificeThis) Description() string {
	return "Sacrifice this permanent"
}

type SacrificeTarget struct {
	cardTypes []string
	target    target.Target
}

func (c SacrificeTarget) TargetSpec() target.TargetSpec {
	return target.PermanentTargetSpec{}
}

func (c SacrificeTarget) isCost() {}

func (c SacrificeTarget) Description() string {
	return fmt.Sprintf("Sacrifice target %v", c.cardTypes)
}

func (c SacrificeTarget) Target() target.Target {
	return c.target
}

func (c SacrificeTarget) WithTarget(target target.Target) CostWithTarget {
	c.target = target
	return c
}
