package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type PassPriorityAction struct {
	playerID string
}

func NewPassPriorityAction(playerID string) PassPriorityAction {
	return PassPriorityAction{
		playerID: playerID,
	}
}

func (a PassPriorityAction) PlayerID() string {
	return a.playerID
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Description() string {
	return "The active player passes priority to the next player."
}

func (a PassPriorityAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Passing priority",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a PassPriorityAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	playerID := game.PriorityPlayerID()
	return []event.GameEvent{event.PassPriorityEvent{
		// TODO: Need to think about how this is managed
		PlayerID: playerID,
	}}, nil
}
