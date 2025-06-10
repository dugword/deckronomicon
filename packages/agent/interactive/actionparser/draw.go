package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type DrawCheatCommand struct {
	Player state.Player
}

func (p *DrawCheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *DrawCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewDrawCheatAction(p.Player), nil
}
