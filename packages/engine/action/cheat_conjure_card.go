package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ConjureCardCheatAction struct {
	player   state.Player
	cardName string
}

func NewConjureCardCheatAction(player state.Player, cardName string) ConjureCardCheatAction {
	return ConjureCardCheatAction{
		player:   player,
		cardName: cardName,
	}
}

func (a ConjureCardCheatAction) Name() string {
	return "Conjure Card"
}

func (a ConjureCardCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.CheatConjureCardEvent{
		PlayerID: a.player.ID(),
		CardName: a.cardName,
	}}, nil
}
