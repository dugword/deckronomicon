package event

import (
	"deckronomicon/packages/game/mtg"
)

// Player Decisions, actual state changes are handled in state change events.
// These events represent intentional choices the player has made during the game.

const (
	EventTypeActivateAbility    = "ActivateAbility"
	EventTypeAssignCombatDamage = "AssignCombatDamage"
	EventTypeCastSpell          = "CastSpell"
	EventTypeConcede            = "Concede"
	EventTypeDeclareAttackers   = "DeclareAttackers"
	EventTypeDeclareBlockers    = "DeclareBlockers"
	EventTypePassPriority       = "PassPriority"
	EventTypePlayLand           = "PlayLand"
)

type PlayerEvent interface{ isPlayerEvent() }

type PlayerBaseEvent struct{}

func (e PlayerBaseEvent) isPlayerEvent() {}

type ActivateAbilityEvent struct {
	PlayerBaseEvent
	PlayerID  string
	AbilityID string
	ObjectID  string
	Zone      mtg.Zone
}

func (e ActivateAbilityEvent) EventType() string {
	return EventTypeActivateAbility
}

type AssignCombatDamageEvent struct {
	PlayerBaseEvent
	PlayerID    string
	Assignments map[string]int // Map of attacker ID to damage assigned
}

func (e AssignCombatDamageEvent) EventType() string {
	return EventTypeAssignCombatDamage
}

type CastSpellEvent struct {
	PlayerBaseEvent
	PlayerID string
	CardID   string
	Zone     mtg.Zone
}

func (e CastSpellEvent) EventType() string {
	return EventTypeCastSpell
}

type ConcedeEvent struct {
	PlayerID string
	PlayerBaseEvent
}

func (e ConcedeEvent) EventType() string {
	return EventTypeConcede
}

type DeclareAttackersEvent struct {
	PlayerBaseEvent
	PlayerID  string
	Attackers []string // List of card IDs that are attacking
}

func (e DeclareAttackersEvent) EventType() string {
	return EventTypeDeclareAttackers
}

type DeclareBlockersEvent struct {
	PlayerBaseEvent
	PlayerID string
	Blockers map[string][]string // Map of attacking card IDs to defending card IDs
}

func (e DeclareBlockersEvent) EventType() string {
	return EventTypeDeclareBlockers
}

type PassPriorityEvent struct {
	PlayerBaseEvent
	PlayerID string
}

func (e PassPriorityEvent) EventType() string {
	return EventTypePassPriority
}

type PlayLandEvent struct {
	PlayerBaseEvent
	PlayerID string
	CardID   string
	Zone     mtg.Zone
}

func (e PlayLandEvent) EventType() string {
	return EventTypePlayLand
}
