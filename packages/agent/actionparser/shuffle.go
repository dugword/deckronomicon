package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ShuffleCheatCommand struct {
	PlayerID string
}

func (p *ShuffleCheatCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *ShuffleCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewShuffleCheatAction(p.PlayerID), nil
}
