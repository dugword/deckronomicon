package event

const (
	EventTypeBeginGame = "BeginGame"
	EventTypeGameOver  = "GameOver"
)

const (
	EventTypeBeginTurn = "BeginTurn"
	EventTypeEndTurn   = "EndTurn"
)

type GameLifecycleEvent interface {
	isGameLifecycleEvent()
}

type GameLifecycleEventBase struct {
}

func (GameLifecycleEventBase) isGameLifecycleEvent() {}

type BeginGameEvent struct {
	GameLifecycleEventBase
}

func (e BeginGameEvent) EventType() string {
	return EventTypeBeginGame
}

func NewBeginGameEvent() BeginGameEvent {
	return BeginGameEvent{}
}

type BeginTurnEvent struct {
	GameLifecycleEventBase
	PlayerID string
}

func (e BeginTurnEvent) EventType() string {
	return EventTypeBeginTurn
}

func NewBeginTurnEvent(playerID string) GameEvent {
	return BeginTurnEvent{
		PlayerID: playerID,
	}
}

type EndTurnEvent struct {
	GameLifecycleEventBase
}

func (e EndTurnEvent) EventType() string {
	return EventTypeEndTurn
}

func NewEndTurnEvent() GameEvent {
	return EndTurnEvent{}
}

type GameOverEvent struct {
	GameLifecycleEventBase
	WinnerID string
}

func (e GameOverEvent) EventType() string {
	return EventTypeGameOver
}

func NewGameOverEvent(winnerID string) GameEvent {
	return GameOverEvent{
		WinnerID: winnerID,
	}
}
