package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type CheatAction struct {
}

func NewCheatAction() CheatAction {
	return CheatAction{}
}

func (a CheatAction) Name() string {
	return "Enable Cheats"
}

func (a CheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{&event.CheatEnabledEvent{
		PlayerID: playerID,
	}}, nil
}
