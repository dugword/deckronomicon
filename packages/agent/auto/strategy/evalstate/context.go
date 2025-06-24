package evalstate

import "deckronomicon/packages/state"

type EvalState struct {
	Game     state.Game
	PlayerID string
	Groups   map[string][]any
	Mode     string
}
