package event

const (
	EventDrawStartingHand      = "DrawStartingHand"
	EventTypePhaseInPhaseOut   = "PhaseInPhaseOut"
	EventTypeCheckDayNight     = "CheckDayNight"
	EventTypeUntapAll          = "UntapAll"
	EventTypeUpkeep            = "Upkeep"
	EventTypeProgressSaga      = "ProgressSaga"
	EventTypeDiscardToHandSize = "DiscardToHandSize"
	EventTypeRemoveDamage      = "RemoveDamage"
)

type TurnBasedActionEvent interface{ isTurnBasedActionEvent() }

type TurnBasedActionEventBase struct{}

func (e *TurnBasedActionEventBase) isTurnBasedActionEvent() {}

type DrawStartingHandEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *DrawStartingHandEvent) EventType() string {
	return EventDrawStartingHand
}

type PhaseInPhaseOutEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *PhaseInPhaseOutEvent) EventType() string {
	return EventTypePhaseInPhaseOut
}

type CheckDayNightEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *CheckDayNightEvent) EventType() string {
	return EventTypeCheckDayNight
}

type UntapAllEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *UntapAllEvent) EventType() string {
	return EventTypeUntapAll
}

type UpkeepEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *UpkeepEvent) EventType() string {
	return EventTypeUpkeep
}

type ProgressSagaEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *ProgressSagaEvent) EventType() string {
	return EventTypeProgressSaga
}

type DiscardToHandSizeEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *DiscardToHandSizeEvent) EventType() string {
	return EventTypeDiscardToHandSize
}

type RemoveDamageEvent struct {
	TurnBasedActionEventBase
	PlayerID string
}

func (e *RemoveDamageEvent) EventType() string {
	return EventTypeRemoveDamage
}
