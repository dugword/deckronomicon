package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type PeekCheatAction struct {
	player state.Player
}

func NewPeekCheatAction(player state.Player) PeekCheatAction {
	return PeekCheatAction{
		player: player,
	}
}

func (a PeekCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a PeekCheatAction) Name() string {
	return "Peek at the top card of your deck"
}

func (a PeekCheatAction) Description() string {
	return "Look at the top card of your deck."
}

func (a PeekCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Peek at the top card of your deck",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a PeekCheatAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Peeked at the top card of the deck",
	}}, nil
}
