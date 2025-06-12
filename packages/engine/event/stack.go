package event

const (
	EventTypeResolveTopObjectOnStack = "ResolveTopObjectOnStack"
)

type StackEvent interface{ isStackEvent() }

type StackBaseEvent struct{}

func (e StackBaseEvent) isStackEvent() {}

type ResolveTopObjectOnStackEvent struct {
	StackBaseEvent
	Name string
	ID   string
}

func (e ResolveTopObjectOnStackEvent) EventType() string {
	return EventTypeResolveTopObjectOnStack
}
