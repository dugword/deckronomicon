package event

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
)

const (
	EventTypeAddMana     = "AddMana"
	EventTypeDiscardCard = "DiscardCard"
	EventTypeDrawCard    = "DrawCard"
	// EventTypeMoveCard                  = "MoveCard"
	EventTypePutCardInHand             = "PutCardInHand"
	EventTypePutCardOnTopOfLibrary     = "PutCardOnTopOfLibrary"
	EventTypePutCardOnBottomOfLibrary  = "PutCardOnBottomOfLibrary"
	EventTypePutPermanentOnBattlefield = "PutPermanentOnBattlefield"
	EventTypePutSpellOnStack           = "PutSpellOnStack"
	EventTypePutSpellInGraveyard       = "PutSpellInGraveyard"
	EventTypePutAbilityOnStack         = "PutAbilityOnStack"
	EventTypeRemoveAbilityFromStack    = "RemoveAbilityFromStack"
	EventTypeResolveManaAbility        = "ResolveManaAbility"
	EventTypeSetActivePlayer           = "SetActivePlayer"
	EventTypeSpendMana                 = "SpendMana"
	EventTypeShuffleDeck               = "ShuffleDeck"
	EventTypeTapPermanent              = "TapPermanent"
	EventTypeUntapPermanent            = "UntapPermanent"
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

type CheatEnabledEvent struct {
	GameStateChangeBaseEvent
	Player string // Player ID who enabled cheats
}

func (e CheatEnabledEvent) EventType() string {
	return EventTypeCheatEnabled
}

/*
type MoveCardEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
	ToZone   mtg.Zone
}

func (e MoveCardEvent) EventType() string {
	return EventTypeMoveCard
}
*/

type PutCardInHandEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutCardInHandEvent) EventType() string {
	return EventTypePutCardInHand
}

type PutCardOnTopOfLibraryEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutCardOnTopOfLibraryEvent) EventType() string {
	return EventTypePutCardOnTopOfLibrary
}

type PutCardOnBottomOfLibraryEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutCardOnBottomOfLibraryEvent) EventType() string {
	return EventTypePutCardOnBottomOfLibrary
}

type PutPermanentOnBattlefieldEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutPermanentOnBattlefieldEvent) EventType() string {
	return EventTypePutPermanentOnBattlefield
}

type PutSpellInGraveyardEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	SpellID  string
}

func (e PutSpellInGraveyardEvent) EventType() string {
	return EventTypePutSpellInGraveyard
}

type PutSpellOnStackEvent struct {
	GameStateChangeBaseEvent
	PlayerID string
	CardID   string
	FromZone mtg.Zone
}

func (e PutSpellOnStackEvent) EventType() string {
	return EventTypePutSpellOnStack
}

type PutAbilityOnStackEvent struct {
	GameStateChangeBaseEvent
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
	GameStateChangeBaseEvent
	PlayerID  string
	AbilityID string
}

func (e RemoveAbilityFromStackEvent) EventType() string {
	return EventTypeRemoveAbilityFromStack
}

type ResolveManaAbilityEvent struct {
	GameStateChangeBaseEvent
	PlayerID    string
	SourceID    string
	AbilityID   string
	FromZone    mtg.Zone
	AbilityName string
	Effects     []definition.EffectSpec
}

func (e ResolveManaAbilityEvent) EventType() string {
	return EventTypeResolveManaAbility
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

type SpendManaEvent struct {
	GameStateChangeBaseEvent
	ManaString string
	PlayerID   string
}

func (e SpendManaEvent) EventType() string {
	return EventTypeSpendMana
}

type TapPermanentEvent struct {
	GameStateChangeBaseEvent
	PlayerID    string
	PermanentID string
}

func (e TapPermanentEvent) EventType() string {
	return EventTypeTapPermanent
}

type UntapPermanentEvent struct {
	GameStateChangeBaseEvent
	PlayerID    string
	PermanentID string
}

func (e UntapPermanentEvent) EventType() string {
	return EventTypeUntapPermanent
}
