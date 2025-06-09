package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type DiscardCheatAction struct {
	player   state.Player
	cardName string
}

func NewDiscardCheatAction(player state.Player, cardName string) DiscardCheatAction {
	return DiscardCheatAction{
		player:   player,
		cardName: cardName,
	}
}

func (a DiscardCheatAction) PlayerID() string {
	return a.player.ID()
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
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Discarded a card from your hand",
	}}, nil
}
