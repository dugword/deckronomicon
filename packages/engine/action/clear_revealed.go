package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

type ClearRevealedAction struct {
	player state.Player
}

func NewClearRevealedAction(player state.Player) ClearRevealedAction {
	return ClearRevealedAction{
		player: player,
	}
}

func (a ClearRevealedAction) PlayerID() string {
	return a.player.ID()
}

func (a ClearRevealedAction) Name() string {
	return "Clear revealed cards"
}

func (a ClearRevealedAction) Description() string {
	return "Clear all revealed cards from your view."
}

func (a ClearRevealedAction) Complete(game state.Game) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.ClearRevealedEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
