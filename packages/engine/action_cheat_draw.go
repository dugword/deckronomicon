package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type DrawCheatAction struct {
	playerID string
}

func NewDrawCheatAction(playerID string) DrawCheatAction {
	return DrawCheatAction{
		playerID: playerID,
	}
}

func (a DrawCheatAction) PlayerID() string {
	return a.playerID
}

func (a DrawCheatAction) Name() string {
	return "Draw a Card"
}

func (a DrawCheatAction) Description() string {
	return "Draw a card from your hand."
}

func (a DrawCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Draw a card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a DrawCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Drew a card from your hand",
	}}, nil
}
