package action

import "deckronomicon/packages/engine/effect"

type ResolutionEnvironment struct {
	EffectRegistry *effect.EffectRegistry
	// definitions maybe should live here too
}
