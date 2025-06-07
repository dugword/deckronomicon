package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ResetLandDropCheatAction struct {
	playerID string
}

func NewResetLandDropCheatAction(playerID string) ResetLandDropCheatAction {
	return ResetLandDropCheatAction{
		playerID: playerID,
	}
}

func (a ResetLandDropCheatAction) PlayerID() string {
	return a.playerID
}

func (a ResetLandDropCheatAction) Name() string {
	return "Reset Land Drop"
}

func (a ResetLandDropCheatAction) Description() string {
	return "Reset the land drop for the turn."
}

func (a ResetLandDropCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Reset Land Drop",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ResetLandDropCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Reset the land drop for the turn",
	}}, nil
}
