package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/state"
)

func ResolveGainLife(
	game *state.Game,
	playerID string,
	gainLife *effect.GainLife,
) (Result, error) {
	events := []event.GameEvent{
		&event.GainLifeEvent{
			PlayerID: playerID,
			Amount:   gainLife.Count,
		},
	}
	return Result{Events: events}, nil
}
