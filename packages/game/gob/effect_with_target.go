package gob

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/target"
)

type EffectWithTarget struct {
	EffectSpec definition.EffectSpec
	Target     target.TargetValue
	SourceID   string
}
