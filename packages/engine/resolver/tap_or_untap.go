package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/state"
)

func ResolveTapOrUntap(
	game state.Game,
	playerID string,
	tapOrUntap effect.TapOrUntap,
	target effect.Target,
) (Result, error) {
	return Result{
		Events: []event.GameEvent{
			event.UntapPermanentEvent{
				PlayerID:    playerID,
				PermanentID: target.ID,
			},
		},
	}, nil
}
