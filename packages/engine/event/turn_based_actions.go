package event

const (
	EventTypePhaseInPhaseOut   = "PhaseInPhaseOut"
	EventTypeCheckDayNight     = "CheckDayNight"
	EventTypeUntapAll          = "UntapAll"
	EventTypeUpkeep            = "Upkeep"
	EventTypeProgressSaga      = "ProgressSaga"
	EventTypeDiscardToHandSize = "DiscardToHandSize"
	EventTypeRemoveDamage      = "RemoveDamage"
)

type TurnBasedActionEvent interface {
	isTurnBasedActionEvent()
}

type TurnBasedActionEventBase struct {
	PlayerID string
}

func (e TurnBasedActionEventBase) isTurnBasedActionEvent() {}

type PhaseInPhaseOutEvent struct {
	TurnBasedActionEventBase
}

func (e PhaseInPhaseOutEvent) EventType() string {
	return EventTypePhaseInPhaseOut
}

func NewPhaseInPhaseOutEvent(playerID string) PhaseInPhaseOutEvent {
	return PhaseInPhaseOutEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type CheckDayNightEvent struct {
	TurnBasedActionEventBase
}

func (e CheckDayNightEvent) EventType() string {
	return EventTypeCheckDayNight
}

func NewCheckDayNightEvent(playerID string) CheckDayNightEvent {
	return CheckDayNightEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type UntapAllEvent struct {
	TurnBasedActionEventBase
}

func (e UntapAllEvent) EventType() string {
	return EventTypeUntapAll
}

func NewUntapAllEvent(playerID string) UntapAllEvent {
	return UntapAllEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type UpkeepEvent struct {
	TurnBasedActionEventBase
}

func (e UpkeepEvent) EventType() string {
	return EventTypeUpkeep
}

func NewUpkeepEvent(playerID string) UpkeepEvent {
	return UpkeepEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type ProgressSagaEvent struct {
	TurnBasedActionEventBase
}

func (e ProgressSagaEvent) EventType() string {
	return EventTypeProgressSaga
}

func NewProgressSagaEvent(playerID string) ProgressSagaEvent {
	return ProgressSagaEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type DiscardToHandSizeEvent struct {
	TurnBasedActionEventBase
}

func (e DiscardToHandSizeEvent) EventType() string {
	return EventTypeDiscardToHandSize
}

func NewDiscardToHandSizeEvent(playerID string) DiscardToHandSizeEvent {
	return DiscardToHandSizeEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}

type RemoveDamageEvent struct {
	TurnBasedActionEventBase
}

func (e RemoveDamageEvent) EventType() string {
	return EventTypeRemoveDamage
}

func NewRemoveDamageEvent(playerID string) RemoveDamageEvent {
	return RemoveDamageEvent{
		TurnBasedActionEventBase: TurnBasedActionEventBase{
			PlayerID: playerID,
		},
	}
}
