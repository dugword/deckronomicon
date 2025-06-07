package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ConcedeCommand struct {
	PlayerID string
}

func (p *ConcedeCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *ConcedeCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewConcedeAction(p.PlayerID), nil
}
