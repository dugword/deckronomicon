package action

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/game/definition"
)

type ResolutionEnvironment struct {
	EffectRegistry *effect.EffectRegistry
	Definitions    map[string]definition.Card
	// definitions maybe should live here too
}
