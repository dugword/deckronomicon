package evalstate

import "deckronomicon/packages/state"

type EvalState struct {
	Game        state.Game
	Player      state.Player
	Definitions map[string][]string
}
