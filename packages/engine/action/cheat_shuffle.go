package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
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

func (a ShuffleCheatAction) Name() string {
	return "Shuffle Deck"
}

func (a ShuffleCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	shuffledCardsIDs := resEnv.RNG.ShuffleCardsIDs(a.player.Library().GetAll())
	return []event.GameEvent{
		event.CheatShuffleDeckEvent{
			PlayerID: a.player.ID(),
		},
		event.ShuffleLibraryEvent{
			PlayerID:         a.player.ID(),
			ShuffledCardsIDs: shuffledCardsIDs,
		},
	}, nil
}
