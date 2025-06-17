package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type CheatAction struct {
	player state.Player
}

func NewCheatAction(player state.Player) CheatAction {
	return CheatAction{
		player: player,
	}
}

func (a CheatAction) PlayerID() string {
	return a.player.ID()
}

func (a CheatAction) Name() string {
	return "Enable Cheats"
}

func (a CheatAction) Description() string {
	return "Enable cheat mode"
}

func (a CheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{event.CheatEnabledEvent{
		Player: a.player.ID(),
	}}, nil
}
