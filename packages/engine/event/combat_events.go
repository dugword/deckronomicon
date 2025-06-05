package event

const (
	EventTypeDeclareAttackers = "DeclareAttackers"
	EventTypeDeclareBlockers  = "DeclareBlockers"
	EventTypeCombatDamage     = "CombatDamage"
)

type CombatEvent interface {
	isCombatEvent()
}

type CombatEventBase struct {
}

func (e CombatEventBase) isCombatEvent() {
}

type DeclareAttackersEvent struct {
	CombatEventBase
}

func (e DeclareAttackersEvent) EventType() string {
	return EventTypeDeclareAttackers
}

type DeclareBlockersEvent struct {
	CombatEventBase
}

func (e DeclareBlockersEvent) EventType() string {
	return EventTypeDeclareBlockers
}

type CombatDamageEvent struct {
	CombatEventBase
}

func (e CombatDamageEvent) EventType() string {
	return EventTypeCombatDamage
}

/*
func NewDeclareAttackersEvent(playerID string, attackers map[string]string) GameEvent {
	// attackerID -> defenderID
	return GameEvent{}
}

func NewDeclareBlockersEvent(playerID string, blockers map[string]string) GameEvent {
	// blockerID -> attackerID
	return GameEvent{}
}

func NewAssignCombatDamageEvent(playerID string, assignments map[string]int) GameEvent {
	// attackerID -> damage
	return GameEvent{}
}
*/
