package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
)

func ResolveRegisterDelayedTriggeredAbility(
	playerID string,
	efct *effect.RegisterDelayedTriggeredAbility,
	source gob.Object,
) (Result, error) {
	events := []event.GameEvent{
		&event.RegisterTriggeredAbilityEvent{
			PlayerID:   playerID,
			SourceName: source.Name(),
			SourceID:   source.ID(),
			Trigger: gob.Trigger{
				EventType: efct.EventType,
				Filter:    gob.Filter(efct.EventFilter),
			},
			OneShot: true,
			Effects: efct.Effects,
		},
	}
	return Result{
		Events: events,
	}, nil
}
