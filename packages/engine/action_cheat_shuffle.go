package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ShuffleCheatAction struct {
	playerID string
}

func NewShuffleCheatAction(playerID string) ShuffleCheatAction {
	return ShuffleCheatAction{
		playerID: playerID,
	}
}

func (a ShuffleCheatAction) PlayerID() string {
	return a.playerID
}

func (a ShuffleCheatAction) Name() string {
	return "Shuffle Deck"
}

func (a ShuffleCheatAction) Description() string {
	return "Shuffle the player's deck."
}

func (a ShuffleCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Shuffle Deck",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ShuffleCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Shuffle the player's deck",
	}}, nil
}
