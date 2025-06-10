package event

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
)

const (
	EventTypeAddMana              = "AddMana"
	EventTypeDiscardCard          = "DiscardCard"
	EventTypeDrawCard             = "DrawCard"
	EventTypeMoveCardToZone       = "MoveCard"
	EventTypePutCardOnBattlefield = "PutCardOnBattlefield"
	EventTypePutCardOnStack       = "PutCardOnStack"
	EventTypeSetActivePlayer      = "SetActivePlayer"
	EventTypeShuffleDeck          = "ShuffleDeck"
	EventTypeTapPermanent         = "TapPermanent"
)

type GameStateChangeEvent interface{ isGameStateChangeEvent() }

type GameStateChangeBaseEvent struct{}

func (e GameStateChangeBaseEvent) isGameStateChangeEvent() {}

type AddManaEvent struct {
	GameStateChangeBaseEvent
	Amount   int
	ManaType mana.ManaType
	PlayerID string
}

func (e AddManaEvent) EventType() string {
	return EventTypeAddMana
}

type DiscardCardEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
}

func (e DiscardCardEvent) EventType() string {
	return EventTypeDiscardCard
}

type DrawCardEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
}

func (e DrawCardEvent) EventType() string {
	return EventTypeDrawCard
}

type MoveCardEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
	ToZone   mtg.Zone
}

func (e MoveCardEvent) EventType() string {
	return EventTypeMoveCardToZone
}

type PutCardOnBattlefieldEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutCardOnBattlefieldEvent) EventType() string {
	return EventTypePutCardOnBattlefield
}

type PutCardOnStackEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutCardOnStackEvent) EventType() string {
	return EventTypePutCardOnStack
}

type SetActivePlayerEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
}

func (e SetActivePlayerEvent) EventType() string {
	return EventTypeSetActivePlayer
}

type ShuffleDeckEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
}

func (e ShuffleDeckEvent) EventType() string {
	return EventTypeShuffleDeck
}

type TapPermanentEvent struct {
	GameStateChangeBaseEvent
	PlayerID    string
	PermanentID string
}

func (e TapPermanentEvent) EventType() string {
	return EventTypeTapPermanent
}
