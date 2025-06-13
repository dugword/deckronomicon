package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
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
	return choose.ChoicePrompt{}, nil
}

func (a PeekCheatAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.CheatPeekEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
