package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type DiscardCheatAction struct {
	playerID string
	cardName string
}

func NewDiscardCheatAction(playerID string, cardName string) DiscardCheatAction {
	return DiscardCheatAction{
		playerID: playerID,
		cardName: cardName,
	}
}

func (a DiscardCheatAction) PlayerID() string {
	return a.playerID
}

func (a DiscardCheatAction) Name() string {
	return "Discard a Card"
}

func (a DiscardCheatAction) Description() string {
	return "Discard a card from your hand."
}

func (a DiscardCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Discard a card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a DiscardCheatAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Discarded a card from your hand",
	}}, nil
}
