package evalstate

import "deckronomicon/packages/state"

type EvalState struct {
	Game        state.Game
	PlayerID    string
	Definitions map[string][]string
	Mode        string
}
