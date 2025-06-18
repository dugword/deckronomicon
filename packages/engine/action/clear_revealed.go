package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ClearRevealedAction struct {
	player state.Player
}

func NewClearRevealedAction(player state.Player) ClearRevealedAction {
	return ClearRevealedAction{
		player: player,
	}
}

func (a ClearRevealedAction) Name() string {
	return "Clear revealed cards"
}

func (a ClearRevealedAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.ClearRevealedEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
