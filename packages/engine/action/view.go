package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type ViewAction struct {
	zone   string
	cardID string
}

func NewViewAction(zone string, cardID string) ViewAction {
	return ViewAction{
		zone:   zone,
		cardID: cardID,
	}
}

func (a ViewAction) Name() string {
	return "View card"
}

func (a ViewAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{&event.NoOpEvent{
		Message: "Viewed card in zone " + a.zone + ": " + a.cardID,
	}}, nil
}
