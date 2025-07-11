package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ClearRevealedAction struct {
}

func NewClearRevealedAction() ClearRevealedAction {
	return ClearRevealedAction{}
}

func (a ClearRevealedAction) Name() string {
	return "Clear revealed cards"
}

func (a ClearRevealedAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{&event.ClearRevealedEvent{
		PlayerID: playerID,
	}}, nil
}
