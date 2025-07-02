package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type LogMessageAction struct {
	message string
}

func NewLogMessageAction(message string) LogMessageAction {
	return LogMessageAction{
		message: message,
	}
}

func (a LogMessageAction) Name() string {
	return "Log a message to the game record"
}

func (a LogMessageAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{
		&event.LogMessageEvent{
			PlayerID: playerID,
			Message:  a.message,
		},
	}, nil
}
