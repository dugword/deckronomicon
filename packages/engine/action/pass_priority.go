package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type PassPriorityAction struct {
	player state.Player
}

func NewPassPriorityAction(player state.Player) PassPriorityAction {
	return PassPriorityAction{
		player: player,
	}
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{event.PassPriorityEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
