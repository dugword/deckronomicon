package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

type ShuffleCheatAction struct {
	player state.Player
}

func NewShuffleCheatAction(player state.Player) ShuffleCheatAction {
	return ShuffleCheatAction{
		player: player,
	}
}

func (a ShuffleCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a ShuffleCheatAction) Name() string {
	return "Shuffle Deck"
}

func (a ShuffleCheatAction) Description() string {
	return "Shuffle the player's deck."
}

func (a ShuffleCheatAction) Complete(game state.Game) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatShuffleDeckEvent{
			PlayerID: a.PlayerID(),
		},
		event.ShuffleDeckEvent{
			PlayerID: a.PlayerID(),
		},
	}, nil
}
