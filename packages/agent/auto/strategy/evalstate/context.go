package evalstate

import "deckronomicon/packages/state"

type EvalState struct {
	Game     *state.Game
	PlayerID string
	Mode     string
}
