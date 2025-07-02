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

func (a PassPriorityAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{&event.PassPriorityEvent{
		PlayerID: playerID,
	}}, nil
}
