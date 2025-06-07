package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type CheatAction struct {
	playerID string
}

func NewCheatAction(playerID string) CheatAction {
	return CheatAction{
		playerID: playerID,
	}
}

func (a CheatAction) PlayerID() string {
	return a.playerID
}

func (a CheatAction) Name() string {
	return "Enable Cheats"
}

func (a CheatAction) Description() string {
	return "Enable cheat mode"
}

func (a CheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Enable Cheats",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a CheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Cheats enabled... you cheater",
	}}, nil
}
