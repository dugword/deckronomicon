package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ConcedeAction struct {
	player state.Player
}

func NewConcedeAction(player state.Player) ConcedeAction {
	return ConcedeAction{
		player: player,
	}
}

func (a ConcedeAction) PlayerID() string {
	return a.player.ID()
}

func (a ConcedeAction) Name() string {
	return "Concede"
}

func (a ConcedeAction) Description() string {
	return "The active player concedes the game."
}

func (a ConcedeAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	opponent, ok := game.GetOpponent(a.player.ID())
	if !ok {
		return nil, fmt.Errorf("opponent for player %q not found", a.player.ID())
	}
	return []event.GameEvent{
		event.ConcedeEvent{PlayerID: a.player.ID()},
		event.EndGameEvent{WinnerID: opponent.ID()},
	}, nil
}
