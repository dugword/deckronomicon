package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
)

func ResolveRegisterDelayedTriggeredAbility(
	playerID string,
	efct *effect.RegisterDelayedTriggeredAbility,
	resolvable state.Resolvable,
) (Result, error) {
	events := []event.GameEvent{
		&event.RegisterTriggeredAbilityEvent{
			PlayerID:   playerID,
			SourceName: resolvable.Name(),
			SourceID:   resolvable.ID(),
			Trigger: gob.Trigger{
				EventType: efct.EventType,
				Filter:    efct.Filter,
			},
			OneShot: true,
			Effects: efct.Effects,
		},
	}
	return Result{
		Events: events,
	}, nil
}
