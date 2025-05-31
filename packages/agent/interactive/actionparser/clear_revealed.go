package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type ClearCommand struct {
	Player state.Player
}

func (p *ClearCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *ClearCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewClearRevealedAction(p.Player), nil
}
