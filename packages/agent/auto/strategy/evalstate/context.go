package evalstate

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type EvalState struct {
	Game            *state.Game
	PlayerID        string
	Mode            string
	MaybeApplyEvent func(game *state.Game, gameEvent event.GameEvent) (*state.Game, error)
}
