package event

// Events go in the event record log so they can be replayed later. That means
// the properties should be serializable and public.

type GameEvent interface {
	EventType() string
}

type Source interface {
	Name() string
}

// TODO: maybe use typed constants for event types
const (
	EventTypeSetNextPlayer = "SetNextPlayer"
)

type SetNextPlayerEvent struct {
}

func (e SetNextPlayerEvent) EventType() string {
	return EventTypeSetNextPlayer
}

type DrawCardEvent struct {
	PlayerID string
	Source   Source
}

func (e DrawCardEvent) EventType() string {
	return EventTypeDrawCard
}

type UntapAllEvent struct {
	PlayerID string
}

func (e UntapAllEvent) EventType() string {
	return EventTypeUntapAll
}

type DeclareAttackersEvent struct {
}

func (e DeclareAttackersEvent) EventType() string {
	return EventTypeDeclareAttackers
}

type DeclareBlockersEvent struct {
}

func (e DeclareBlockersEvent) EventType() string {
	return EventTypeDeclareBlockers
}

type CombatDamageEvent struct {
}

func (e CombatDamageEvent) EventType() string {
	return EventTypeCombatDamage
}

type CleanupStepEvent struct {
}

func (e CleanupStepEvent) EventType() string {
	return EventTypeCleanupStep
}
