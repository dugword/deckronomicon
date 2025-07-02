package resenv

import (
	"deckronomicon/packages/engine/rng"
	"deckronomicon/packages/game/definition"
)

// This is a good pattern that I think I like, finish implementing this.
// Then figure out how to pass effect registry... maybe.
// Do actions need effect registry?mm

type ResEnv struct {
	RNG         *rng.RNG
	Definitions map[string]*definition.Card
	//effectRegistry *effect.EffectRegistry
}
