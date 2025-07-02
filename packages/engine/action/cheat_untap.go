package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type UntapCheatAction struct {
	permanentID string
}

func NewUntapCheatAction(permanentID string) UntapCheatAction {
	return UntapCheatAction{
		permanentID: permanentID,
	}
}

func (a UntapCheatAction) Name() string {
	return "Untap target permanent"
}

func (a UntapCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		&event.CheatUntapEvent{
			PlayerID: playerID,
		},
		&event.UntapPermanentEvent{
			PlayerID:    playerID,
			PermanentID: a.permanentID,
		},
	}, nil
}
