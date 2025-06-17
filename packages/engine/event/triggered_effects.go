package event

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

const (
	EventTypeRegisterTriggeredEffect = "RegisterTriggeredEffect"
	EventTypeRemoveTriggeredEffect   = "RemoveTriggeredEffect"
)

type TriggeredEffectEvent interface{ isTriggeredEffectEvent() }

type TriggeredEffectBaseEvent struct{}

func (e TriggeredEffectBaseEvent) isTriggeredEffectEvent() {}

type RegisterTriggeredEffectEvent struct {
	TriggeredEffectBaseEvent
	PlayerID    string
	SourceName  string
	SourceID    string
	Trigger     state.Trigger
	EffectSpecs []definition.EffectSpec
	Duration    mtg.Duration
	OneShot     bool
}

func (e RegisterTriggeredEffectEvent) EventType() string {
	return EventTypeRegisterTriggeredEffect
}

type RemoveTriggeredEffectEvent struct {
	TriggeredEffectBaseEvent
	ID string
}

func (e RemoveTriggeredEffectEvent) EventType() string {
	return EventTypeRemoveTriggeredEffect
}
