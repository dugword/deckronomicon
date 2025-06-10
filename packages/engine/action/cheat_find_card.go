package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type FindCardCheatAction struct {
	player   state.Player
	cardName string
}

func NewFindCardCheatAction(player state.Player, cardName string) FindCardCheatAction {
	return FindCardCheatAction{
		player:   player,
		cardName: cardName,
	}
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) PlayerID() string {
	return a.player.ID()
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
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Found a card from library and added it to your hand",
	}}, nil
}
