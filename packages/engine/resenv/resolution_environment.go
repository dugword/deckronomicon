package resenv

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/rng"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/state"
)

// This is a good pattern that I think I like, finish implementing this.
// Then figure out how to pass effect registry... maybe.
// Do actions need effect registry?mm

type MaybeApplyEvent func(game *state.Game, gameEvent event.GameEvent) (*state.Game, error)

type ResEnv struct {
	RNG             *rng.RNG
	Definitions     map[string]*definition.Card
	MaybeApplyEvent MaybeApplyEvent
}
