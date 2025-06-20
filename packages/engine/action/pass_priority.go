package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type PassPriorityAction struct {
	playerID string
}

func NewPassPriorityAction(playerID string) PassPriorityAction {
	return PassPriorityAction{
		playerID: playerID,
	}
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{event.PassPriorityEvent{
		PlayerID: a.playerID,
	}}, nil
}
