package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type ResetLandDropCommand struct {
	playerID string
}

func (p *ResetLandDropCommand) Build(game *state.Game, playerID string) (engine.Action, error) {
	return action.NewResetLandDropCheatAction(), nil
}
