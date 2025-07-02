package event

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
)

const (
	EventTypeResolveTopObjectOnStack       = "ResolveTopObjectOnStack"
	EventTypePutSpellOnStack               = "PutSpellOnStack"
	EventTypePutCopiedSpellOnStack         = "PutCopiedSpellOnStack"
	EventTypePutAbilityOnStack             = "PutAbilityOnStack"
	EventTypeRemoveSpellOrAbilityFromStack = "RemoveSpellOrAbilityFromStack"
	EventTypeSpellOrAbilityFizzles         = "SpellOrAbilityFizzles"
)

type StackEvent interface{ isStackEvent() }

type StackBaseEvent struct{}

func (e *StackBaseEvent) isStackEvent() {}

type ResolveTopObjectOnStackEvent struct {
	StackBaseEvent
	Name string
	ID   string
}

func (e *ResolveTopObjectOnStackEvent) EventType() string {
	return EventTypeResolveTopObjectOnStack
}

type PutSpellInExileEvent struct {
	StackBaseEvent
	PlayerID string
	SpellID  string
}

type PutSpellInGraveyardEvent struct {
	StackBaseEvent
	PlayerID string
	SpellID  string
}

type PutCopiedSpellOnStackEvent struct {
	StackBaseEvent
	PlayerID          string
	SpellID           string
	FromZone          mtg.Zone
	EffectWithTargets []*effect.EffectWithTarget
}

func (e *PutCopiedSpellOnStackEvent) EventType() string {
	return EventTypePutCopiedSpellOnStack
}

type PutSpellOnStackEvent struct {
	StackBaseEvent
	PlayerID          string
	CardID            string
	FromZone          mtg.Zone
	EffectWithTargets []*effect.EffectWithTarget
	Flashback         bool
}

func (e *PutSpellOnStackEvent) EventType() string {
	return EventTypePutSpellOnStack
}

type PutAbilityOnStackEvent struct {
	StackBaseEvent
	PlayerID          string
	SourceID          string
	AbilityID         string
	FromZone          mtg.Zone
	AbilityName       string
	EffectWithTargets []*effect.EffectWithTarget
}

func (e *PutAbilityOnStackEvent) EventType() string {
	return EventTypePutAbilityOnStack
}

type RemoveSpellOrAbilityFromStackEvent struct {
	StackBaseEvent
	PlayerID string
	ObjectID string
}

func (e *RemoveSpellOrAbilityFromStackEvent) EventType() string {
	return EventTypeRemoveSpellOrAbilityFromStack
}

type SpellOrAbilityFizzlesEvent struct {
	StackBaseEvent
	PlayerID string
	ObjectID string
}

func (e *SpellOrAbilityFizzlesEvent) EventType() string {
	return EventTypeSpellOrAbilityFizzles
}
