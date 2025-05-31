package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

type UntapCheatAction struct {
	player    state.Player
	permanent gob.Permanent
}

func NewUntapCheatAction(player state.Player, permanent gob.Permanent) UntapCheatAction {
	return UntapCheatAction{
		player:    player,
		permanent: permanent,
	}
}

func (a UntapCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a UntapCheatAction) Name() string {
	return "Untap target permanent"
}

func (a UntapCheatAction) Description() string {
	return "Untap target permanent."
}
func (a UntapCheatAction) Complete(game state.Game) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatUntapEvent{
			PlayerID: a.player.ID(),
		},
		event.UntapPermanentEvent{
			PlayerID:    a.player.ID(),
			PermanentID: a.permanent.ID(),
		},
	}, nil
}
