package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type FindCardCheatAction struct {
	cardID string
}

func NewFindCardCheatAction(cardID string) FindCardCheatAction {
	return FindCardCheatAction{
		cardID: cardID,
	}
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	player := game.GetPlayer(playerID)
	card, ok := player.Library().Get(a.cardID)
	if !ok {
		return nil, fmt.Errorf("card %q not found in library", a.cardID)
	}
	return []event.GameEvent{
		&event.CheatFindCardEvent{
			PlayerID: playerID,
			CardID:   a.cardID,
		},
		&event.PutCardInHandEvent{
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
			PlayerID: playerID,
		},
	}, nil
}
