package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ResetLandDropCommand struct {
	Player state.Player
}

func (p *ResetLandDropCommand) IsComplete() bool {
	return p.Player.ID() != ""
}
func (p *ResetLandDropCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return engine.NewResetLandDropCheatAction(player), nil
}
