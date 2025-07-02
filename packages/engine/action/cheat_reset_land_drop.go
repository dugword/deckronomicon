package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ResetLandDropCheatAction struct {
}

func NewResetLandDropCheatAction() ResetLandDropCheatAction {
	return ResetLandDropCheatAction{}
}

func (a ResetLandDropCheatAction) Name() string {
	return "Reset Land Drop"
}

func (a ResetLandDropCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{&event.CheatResetLandDropEvent{
		PlayerID: playerID,
	}}, nil
}
