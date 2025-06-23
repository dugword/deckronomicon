package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

type UntapCheatAction struct {
	permanent gob.Permanent
}

func NewUntapCheatAction(permanent gob.Permanent) UntapCheatAction {
	return UntapCheatAction{
		permanent: permanent,
	}
}

func (a UntapCheatAction) Name() string {
	return "Untap target permanent"
}

func (a UntapCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatUntapEvent{
			PlayerID: player.ID(),
		},
		event.UntapPermanentEvent{
			PlayerID:    player.ID(),
			PermanentID: a.permanent.ID(),
		},
	}, nil
}
