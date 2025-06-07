package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type PassPriorityCommand struct {
	PlayerID string
}

func (p *PassPriorityCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *PassPriorityCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewPassPriorityAction(p.PlayerID), nil
}
