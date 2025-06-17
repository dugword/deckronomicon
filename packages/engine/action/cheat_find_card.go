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
	player state.Player
	card   gob.Card
}

func NewFindCardCheatAction(player state.Player, cardInZone gob.Card) FindCardCheatAction {
	return FindCardCheatAction{
		player: player,
		card:   cardInZone,
	}
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a FindCardCheatAction) Description() string {
	return "Find a card into your hand."
}

func (a FindCardCheatAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	card, ok := a.player.Library().Get(a.card.ID())
	if !ok {
		return nil, fmt.Errorf("card %q not found in library", a.card.ID())
	}
	return []event.GameEvent{
		event.CheatFindCardEvent{
			PlayerID: a.player.ID(),
			CardID:   a.card.ID(),
		},
		event.PutCardInHandEvent{
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
			PlayerID: a.player.ID(),
		},
	}, nil
}
