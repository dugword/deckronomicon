package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type Action interface {
	Name() string
	Description() string
	GetPrompt(state.Game) (choose.ChoicePrompt, error)
	Complete(state.Game, []choose.Choice) ([]event.GameEvent, error)
	PlayerID() string
}

type PassAction struct {
	playerID string
}

func NewPassAction(playerID string) PassAction {
	return PassAction{
		playerID: playerID,
	}
}

func (a PassAction) PlayerID() string {
	return a.playerID
}

func (a PassAction) Name() string {
	return "Pass Priority"
}

func (a PassAction) Description() string {
	return "The active player passes priority to the next player."
}

func (a PassAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Passing priority",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a PassAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	playerID := game.PriorityPlayerID()
	return []event.GameEvent{event.PassPriorityEvent{
		// TODO: Need to think about how this is managed
		PlayerID: playerID,
	}}, nil
}
