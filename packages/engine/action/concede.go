package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ConcedeAction struct {
	player state.Player
}

func NewConcedeAction() ConcedeAction {
	return ConcedeAction{}
}

func (a ConcedeAction) Name() string {
	return "Concede"
}

// TODO: This only works for 2 player games.
func (a ConcedeAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	opponent := game.GetOpponent(player.ID())
	return []event.GameEvent{
		event.ConcedeEvent{PlayerID: a.player.ID()},
		event.EndGameEvent{WinnerID: opponent.ID()},
	}, nil
}
