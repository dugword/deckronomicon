package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type CheatCommand struct {
	Player state.Player
}

func (p *CheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *CheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewCheatAction(p.Player), nil
}
