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

func (a LogMessageAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{
		event.LogMessageEvent{
			PlayerID: player.ID(),
			Message:  a.message,
		},
	}, nil
}
