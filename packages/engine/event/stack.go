package event

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
)

const (
	EventTypeResolveTopObjectOnStack = "ResolveTopObjectOnStack"
	EventTypePutSpellOnStack         = "PutSpellOnStack"
	EventTypePutSpellInGraveyard     = "PutSpellInGraveyard"
	EventTypePutSpellInExile         = "PutSpellInExile"
	EventTypePutAbilityOnStack       = "PutAbilityOnStack"
	EventTypeRemoveAbilityFromStack  = "RemoveAbilityFromStack"
)

type StackEvent interface{ isStackEvent() }

type StackBaseEvent struct{}

func (e StackBaseEvent) isStackEvent() {}

type ResolveTopObjectOnStackEvent struct {
	StackBaseEvent
	Name string
	ID   string
}

func (e ResolveTopObjectOnStackEvent) EventType() string {
	return EventTypeResolveTopObjectOnStack
}

type PutSpellInExileEvent struct {
	StackBaseEvent
	PlayerID string
	SpellID  string
}

func (e PutSpellInExileEvent) EventType() string {
	return EventTypePutSpellInExile
}

type PutSpellInGraveyardEvent struct {
	StackBaseEvent
	PlayerID string
	SpellID  string
}

func (e PutSpellInGraveyardEvent) EventType() string {
	return EventTypePutSpellInGraveyard
}

type PutSpellOnStackEvent struct {
	StackBaseEvent
	PlayerID  string
	CardID    string
	FromZone  mtg.Zone
	Flashback bool
}

func (e PutSpellOnStackEvent) EventType() string {
	return EventTypePutSpellOnStack
}

type PutAbilityOnStackEvent struct {
	StackBaseEvent
	PlayerID    string
	SourceID    string
	AbilityID   string
	FromZone    mtg.Zone
	AbilityName string
	Effects     []definition.EffectSpec
}

func (e PutAbilityOnStackEvent) EventType() string {
	return EventTypePutAbilityOnStack
}

type RemoveAbilityFromStackEvent struct {
	StackBaseEvent
	PlayerID  string
	AbilityID string
}

func (e RemoveAbilityFromStackEvent) EventType() string {
	return EventTypeRemoveAbilityFromStack
}
