package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type PeekCheatAction struct {
}

func NewPeekCheatAction(playerID string) PeekCheatAction {
	return PeekCheatAction{}
}

func (a PeekCheatAction) Name() string {
	return "Peek at the top card of your deck"
}

func (a PeekCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{&event.CheatPeekEvent{
		PlayerID: playerID,
	}}, nil
}
