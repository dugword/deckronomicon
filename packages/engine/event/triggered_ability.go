package event

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
)

const (
	EventTypeRegisterTriggeredAbility = "RegisterTriggeredAbility"
	EventTypeRemoveTriggeredAbility   = "RemoveTriggeredAbility"
)

type TriggeredAbilityEvent interface{ isTriggeredAbilityEvent() }

type TriggeredAbilityBaseEvent struct{}

func (e TriggeredAbilityBaseEvent) isTriggeredAbilityEvent() {}

type RegisterTriggeredAbilityEvent struct {
	TriggeredAbilityBaseEvent
	PlayerID   string
	SourceName string
	SourceID   string
	Trigger    gob.Trigger
	Effects    []effect.Effect
	Duration   mtg.Duration
	OneShot    bool
}

func (e RegisterTriggeredAbilityEvent) EventType() string {
	return EventTypeRegisterTriggeredAbility
}

type RemoveTriggeredAbilityEvent struct {
	TriggeredAbilityBaseEvent
	ID string
}

func (e RemoveTriggeredAbilityEvent) EventType() string {
	return EventTypeRemoveTriggeredAbility
}
