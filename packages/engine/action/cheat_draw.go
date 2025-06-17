package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type DrawCheatAction struct {
	player state.Player
}

func NewDrawCheatAction(player state.Player) DrawCheatAction {
	return DrawCheatAction{
		player: player,
	}
}

func (a DrawCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a DrawCheatAction) Name() string {
	return "Draw a Card"
}

func (a DrawCheatAction) Description() string {
	return "Draw a card from your hand."
}

func (a DrawCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatDrawEvent{
			PlayerID: a.player.ID(),
		},
		event.DrawCardEvent{
			PlayerID: a.player.ID(),
		},
	}, nil
}
