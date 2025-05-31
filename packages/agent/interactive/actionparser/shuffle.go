package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type ShuffleCheatCommand struct {
	Player state.Player
}

func (p *ShuffleCheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *ShuffleCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewShuffleCheatAction(p.Player), nil
}
