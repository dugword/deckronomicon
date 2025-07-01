package actionparser

import (
	"deckronomicon/packages/engine/action"
	"fmt"
)

func parseLogMessage(
	message string,
) (action.LogMessageAction, error) {
	if message == "" {
		return action.LogMessageAction{}, fmt.Errorf("log message command requires a message to log")
	}
	return action.NewLogMessageAction(message), nil
}
