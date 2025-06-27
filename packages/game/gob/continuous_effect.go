package gob

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
)

// tbd where this goes, I think it's more like the stack
// and should live in state
type ContinuousEffect struct {
	ID       string
	Duration mtg.Duration
	// Layer mtg.Layer
	Source Object
	// Targets []target.TargetSpec
	Effect []effect.Effect
	// Timestamp int (these need to be applied in order)
}
