package event

const (
	EventTypeAllPlayersPassedPriority = "AllPlayersPassedPriority"
	EventTypePassPriority             = "PassPriority"
	EventTypeReceivePriorityEvent     = "ReceivePriorityEvent"
	EventTypeResetPriorityPasses      = "ResetPriorityPasses"
)

type PriorityEvent interface {
	isPriorityEvent()
}

type priorityEvent struct{}

func (e priorityEvent) isPriorityEvent() {}

type AllPlayersPassedPriorityEvent struct {
	priorityEvent
}

func (e AllPlayersPassedPriorityEvent) EventType() string {
	return EventTypeAllPlayersPassedPriority
}

type PassPriorityEvent struct {
	priorityEvent
	PlayerID string
}

func (e PassPriorityEvent) EventType() string {
	return EventTypePassPriority
}

type ReceivePriorityEvent struct {
	priorityEvent
	PlayerID string
}

func (e ReceivePriorityEvent) EventType() string {
	return EventTypeReceivePriorityEvent
}

type ResetPriorityPassesEvent struct {
	priorityEvent
}

func (e ResetPriorityPassesEvent) EventType() string {
	return EventTypeResetPriorityPasses
}
