package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

type DiscardCheatAction struct {
	player state.Player
	card   gob.Card
}

func NewDiscardCheatAction(player state.Player, card gob.Card) DiscardCheatAction {
	return DiscardCheatAction{
		player: player,
		card:   card,
	}
}

func (a DiscardCheatAction) Name() string {
	return "Discard a Card"
}

func (a DiscardCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatDiscardEvent{
			PlayerID: a.player.ID(),
		},
		event.DiscardCardEvent{
			PlayerID: a.player.ID(),
			CardID:   a.card.ID(),
		},
	}, nil
}
