package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type AddManaCheatAction struct {
	playerID string
	mana     string
}

func NewAddManaCheatAction(playerID string, mana string) AddManaCheatAction {
	return AddManaCheatAction{
		playerID: playerID,
		mana:     mana,
	}
}

func (a AddManaCheatAction) PlayerID() string {
	return a.playerID
}

func (a AddManaCheatAction) Name() string {
	return "CHEAT: Add Mana"
}

func (a AddManaCheatAction) Description() string {
	return "CHEAT: The active player adds mana."
}

func (a AddManaCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "CHEAT: Adding mana",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a AddManaCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "CHEAT: Added mana",
	}}, nil
}
