package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
	"fmt"
)

type ConjureCardCheatAction struct {
	cardName string
}

func NewConjureCardCheatAction(cardName string) ConjureCardCheatAction {
	return ConjureCardCheatAction{
		cardName: cardName,
	}
}

func (a ConjureCardCheatAction) Name() string {
	return "Conjure Card"
}

func (a ConjureCardCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.CheatConjureCardEvent{
		PlayerID: player.ID(),
		CardName: a.cardName,
	}}, nil
}
