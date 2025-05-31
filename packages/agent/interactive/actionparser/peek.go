package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type PeekCheatCommand struct {
	Player state.Player
}

func (p *PeekCheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *PeekCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewPeekCheatAction(player), nil
}
