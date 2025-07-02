package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type DrawCheatAction struct {
}

func NewDrawCheatAction() DrawCheatAction {
	return DrawCheatAction{}
}

func (a DrawCheatAction) Name() string {
	return "Draw a Card"
}

func (a DrawCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		&event.CheatDrawEvent{
			PlayerID: playerID,
		},
		&event.DrawCardEvent{
			PlayerID: playerID,
		},
	}, nil
}
