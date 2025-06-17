package state

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

type ContinuousEffect struct {
	ID       string
	Duration mtg.Duration
	// Layer mtg.Layer
	Source query.Object
	// Targets []target.TargetSpec
	EffectSpecs []definition.EffectSpec
	// Timestamp int (these need to be applied in order)
}
