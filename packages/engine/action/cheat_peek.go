package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
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

func (a PeekCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.CheatPeekEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
