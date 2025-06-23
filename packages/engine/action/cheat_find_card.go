package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type FindCardCheatAction struct {
	card gob.Card
}

func NewFindCardCheatAction(card gob.Card) FindCardCheatAction {
	return FindCardCheatAction{
		card: card,
	}
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	card, ok := player.Library().Get(a.card.ID())
	if !ok {
		return nil, fmt.Errorf("card %q not found in library", a.card.ID())
	}
	return []event.GameEvent{
		event.CheatFindCardEvent{
			PlayerID: player.ID(),
			CardID:   a.card.ID(),
		},
		event.PutCardInHandEvent{
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
			PlayerID: player.ID(),
		},
	}, nil
}
