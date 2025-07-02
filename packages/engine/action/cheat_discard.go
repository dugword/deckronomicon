package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type DiscardCheatAction struct {
	cardID string
}

func NewDiscardCheatAction(cardID string) DiscardCheatAction {
	return DiscardCheatAction{
		cardID: cardID,
	}
}

func (a DiscardCheatAction) Name() string {
	return "Discard a Card"
}

func (a DiscardCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		&event.CheatDiscardEvent{
			PlayerID: playerID,
		},
		&event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   a.cardID,
		},
	}, nil
}
