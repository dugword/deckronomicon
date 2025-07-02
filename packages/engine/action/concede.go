package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ConcedeAction struct {
	playerID string
}

func NewConcedeAction() ConcedeAction {
	return ConcedeAction{}
}

func (a ConcedeAction) Name() string {
	return "Concede"
}

// TODO: This only works for 2 player games.
func (a ConcedeAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	opponentID := game.GetOpponent(playerID).ID()
	return []event.GameEvent{
		&event.ConcedeEvent{PlayerID: a.playerID},
		&event.EndGameEvent{WinnerID: opponentID},
	}, nil
}
