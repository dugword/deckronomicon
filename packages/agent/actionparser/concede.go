package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ConcedeCommand struct {
	Player state.Player
}

func (p *ConcedeCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *ConcedeCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return engine.NewConcedeAction(p.Player), nil
}
