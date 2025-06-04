package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type Action interface {
	Name() string
	Description() string
	Complete(state.Game, choose.Choice) (event.GameEvent, error)
}

type PassAction struct {
}

func (a PassAction) Name() string {
	return "Pass Priority"
}

func (a PassAction) Description() string {
	return "The active player passes priority to the next player."
}

func (a PassAction) GetPrompt(state state.Game, player state.Player) choose.ChoicePrompt {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Passing priority",
		Choices:  nil,
		Optional: false,
	}
}
func (a PassAction) Complete(game state.Game, choice choose.Choice) (event.GameEvent, error) {
	playerID := game.PriorityPlayerID()
	return event.PassPriorityEvent{
		// TODO: Need to think about how this is managed
		PlayerID: playerID,
	}, nil
}
