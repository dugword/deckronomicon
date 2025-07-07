package event

import "deckronomicon/packages/game/mtg"

const (
	EventTypeEnteredTheBattlefield = "EnteredTheBattlefield"
	EventTypeLeftTheBattlefield    = "LeftTheBattlefield"
	EventTypeTapped                = "PermanentTapped"
	EventTypeUntapped              = "PermanentUntapped"
	EventTypeDeath                 = "Death"
	EventTypeLandTappedForMana     = "LandTappedForMana"
)

type TriggerEvent interface {
	isTriggerEvent()
}

type TriggerBaseEvent struct{}

func (e TriggerBaseEvent) isTriggerEvent() {}

type EnteredTheBattlefieldEvent struct {
	TriggerBaseEvent
	ControllerID string
	PermanentID  string
	CardTypes    []mtg.CardType
	Subtypes     []mtg.Subtype
	Supertypes   []mtg.Supertype
}

func (e EnteredTheBattlefieldEvent) EventType() string {
	return EventTypeEnteredTheBattlefield
}

type LeftTheBattlefieldEvent struct {
	TriggerBaseEvent
	ControllerID string
	PermanentID  string
	CardTypes    []mtg.CardType
	Subtypes     []mtg.Subtype
	Supertypes   []mtg.Supertype
}

func (e LeftTheBattlefieldEvent) EventType() string {
	return EventTypeLeftTheBattlefield
}

type TappedEvent struct {
	TriggerBaseEvent
	ControllerID string
	PermanentID  string
	CardTypes    []mtg.CardType
	Subtypes     []mtg.Subtype
	Supertypes   []mtg.Supertype
}

func (e TappedEvent) EventType() string {
	return EventTypeTapped
}

type UntappedEvent struct {
	TriggerBaseEvent
	ControllerID string
	PermanentID  string
	CardTypes    []mtg.CardType
	Subtypes     []mtg.Subtype
	Supertypes   []mtg.Supertype
}

func (e UntappedEvent) EventType() string {
	return EventTypeUntapped
}

type DeathEvent struct {
	TriggerBaseEvent
	ControllerID string
	OwnerID      string
	PermanentID  string
	CardID       string
	CardTypes    []mtg.CardType
	Subtypes     []mtg.Subtype
	Supertypes   []mtg.Supertype
}

func (e DeathEvent) EventType() string {
	return EventTypeDeath
}

type LandTappedForManaEvent struct {
	TriggerBaseEvent
	PlayerID string
	ObjectID string
	Subtypes []mtg.Subtype
}

func (e *LandTappedForManaEvent) EventType() string {
	return EventTypeLandTappedForMana
}
