package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

type DiscardCheatAction struct {
	card gob.Card
}

func NewDiscardCheatAction(card gob.Card) DiscardCheatAction {
	return DiscardCheatAction{
		card: card,
	}
}

func (a DiscardCheatAction) Name() string {
	return "Discard a Card"
}

func (a DiscardCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatDiscardEvent{
			PlayerID: player.ID(),
		},
		event.DiscardCardEvent{
			PlayerID: player.ID(),
			CardID:   a.card.ID(),
		},
	}, nil
}
