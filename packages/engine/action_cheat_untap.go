package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type UntapCheatAction struct {
	playerID    string
	permanentID string
}

func NewUntapCheatAction(playerID string, permanentID string) UntapCheatAction {
	return UntapCheatAction{
		playerID:    playerID,
		permanentID: permanentID,
	}
}

func (a UntapCheatAction) PlayerID() string {
	return a.playerID
}

func (a UntapCheatAction) Name() string {
	return "Untap target permanent"
}

func (a UntapCheatAction) Description() string {
	return "Untap target permanent."
}

func (a UntapCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Untap Permanent",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a UntapCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Untap target permanent",
	}}, nil
}
