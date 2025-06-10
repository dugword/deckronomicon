package event

const (
	EventTypeAllPlayersPassedPriority = "AllPlayersPassedPriority"
	EventTypeReceivePriorityEvent     = "ReceivePriorityEvent"
	EventTypeResetPriorityPasses      = "ResetPriorityPasses"
)

type PriorityEvent interface{ isPriorityEvent() }

type PriorityBaseEvent struct{}

func (e PriorityBaseEvent) isPriorityEvent() {}

type AllPlayersPassedPriorityEvent struct {
	PriorityBaseEvent
}

func (e AllPlayersPassedPriorityEvent) EventType() string {
	return EventTypeAllPlayersPassedPriority
}

type ReceivePriorityEvent struct {
	PriorityBaseEvent
	PlayerID string
}

func (e ReceivePriorityEvent) EventType() string {
	return EventTypeReceivePriorityEvent
}

type ResetPriorityPassesEvent struct {
	PriorityBaseEvent
}

func (e ResetPriorityPassesEvent) EventType() string {
	return EventTypeResetPriorityPasses
}
