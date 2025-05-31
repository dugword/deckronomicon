package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type PassPriorityCommand struct {
	Player state.Player
}

func (p *PassPriorityCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *PassPriorityCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewPassPriorityAction(p.Player), nil
}
