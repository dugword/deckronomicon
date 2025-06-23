package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type PassPriorityAction struct{}

func NewPassPriorityAction() PassPriorityAction {
	return PassPriorityAction{}
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{event.PassPriorityEvent{
		PlayerID: player.ID(),
	}}, nil
}
