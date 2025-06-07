package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type FindCardCheatAction struct {
	playerID string
	cardName string
}

func NewFindCardCheatAction(playerID string, cardName string) FindCardCheatAction {
	return FindCardCheatAction{
		playerID: playerID,
		cardName: cardName,
	}
}

func (a FindCardCheatAction) PlayerID() string {
	return a.playerID
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) Description() string {
	return "Find a card into your hand."
}

func (a FindCardCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Find a card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a FindCardCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Found a card from library and added it to your hand",
	}}, nil
}
