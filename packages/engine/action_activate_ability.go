package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ActivateAbilityAction struct {
	playerID string
}

func NewActivateAbilityAction(playerID string) ActivateAbilityAction {
	return ActivateAbilityAction{
		playerID: playerID,
	}
}

func (a ActivateAbilityAction) PlayerID() string {
	return a.playerID
}

func (a ActivateAbilityAction) Name() string {
	return "Activate Ability"
}

func (a ActivateAbilityAction) Description() string {
	return "The active player activates an ability."
}

func (a ActivateAbilityAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Activating an ability",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ActivateAbilityAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Activated an ability",
	}}, nil
}
