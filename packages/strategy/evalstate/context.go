package evalstate

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/player"
)

type EvalState struct {
	State       *engine.GameState
	Player      *player.Player
	Definitions map[string][]string
}
