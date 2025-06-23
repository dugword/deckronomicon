package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ShuffleCheatAction struct {
}

func NewShuffleCheatAction(player state.Player) ShuffleCheatAction {
	return ShuffleCheatAction{}
}

func (a ShuffleCheatAction) Name() string {
	return "Shuffle Deck"
}

func (a ShuffleCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	var cardIDs []string
	for _, card := range player.Library().GetAll() {
		cardIDs = append(cardIDs, card.ID())
	}
	shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
	return []event.GameEvent{
		event.CheatShuffleDeckEvent{
			PlayerID: player.ID(),
		},
		event.ShuffleLibraryEvent{
			PlayerID:         player.ID(),
			ShuffledCardsIDs: shuffledCardsIDs,
		},
	}, nil
}
