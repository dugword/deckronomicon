package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ResetLandDropCommand struct {
	PlayerID string
	Card     string
}

func (p *ResetLandDropCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Card != ""
}
func (p *ResetLandDropCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewResetLandDropCheatAction(p.PlayerID), nil
}
