package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ConcedeAction struct {
	player state.Player
}

func NewConcedeAction(player state.Player) ConcedeAction {
	return ConcedeAction{
		player: player,
	}
}

func (a ConcedeAction) Name() string {
	return "Concede"
}

// TODO: This only works for 2 player games.
func (a ConcedeAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	opponent := game.GetOpponent(a.player.ID())
	return []event.GameEvent{
		event.ConcedeEvent{PlayerID: a.player.ID()},
		event.EndGameEvent{WinnerID: opponent.ID()},
	}, nil
}
