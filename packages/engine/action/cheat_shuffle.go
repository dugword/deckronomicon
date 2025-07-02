package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ShuffleCheatAction struct {
}

func NewShuffleCheatAction(playerID string) ShuffleCheatAction {
	return ShuffleCheatAction{}
}

func (a ShuffleCheatAction) Name() string {
	return "Shuffle Deck"
}

func (a ShuffleCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	player := game.GetPlayer(playerID)
	var cardIDs []string
	for _, card := range player.Library().GetAll() {
		cardIDs = append(cardIDs, card.ID())
	}
	shuffledCardsIDs := resEnv.RNG.ShuffleIDs(cardIDs)
	return []event.GameEvent{
		&event.CheatShuffleDeckEvent{
			PlayerID: playerID,
		},
		&event.ShuffleLibraryEvent{
			PlayerID:         playerID,
			ShuffledCardsIDs: shuffledCardsIDs,
		},
	}, nil
}
