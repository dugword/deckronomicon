package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ViewAction struct {
	player state.Player
	zone   string
	cardID string
}

func NewViewAction(player state.Player, zone string, cardID string) ViewAction {
	return ViewAction{
		player: player,
		zone:   zone,
		cardID: cardID,
	}
}

func (a ViewAction) Name() string {
	return "View card"
}

func (a ViewAction) Complete(game state.Game, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Viewed card in zone " + a.zone + ": " + a.cardID,
	}}, nil
}
