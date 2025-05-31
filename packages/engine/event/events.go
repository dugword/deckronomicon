package event

// Events go in the event record log so they can be replayed later. That means
// the properties should be serializable and public.

const (
	EventTypeNoOp      = "NoOp"
	EventTypeMilestone = "Milestone"
)

type GameEvent interface {
	EventType() string
}

// NoOpEvent is a placeholder event that does nothing. It is used during
// development when a place holder event is needed.

type NoOpEvent struct {
	Message string
}

func (e NoOpEvent) EventType() string {
	return "NoOp"
}

// MilestoneEvent is used to mark significant points in the game as defined by the player,
// this is used for analytics or tracking purposes.
type MilestoneEvent struct {
	Message string
}

func (e MilestoneEvent) EventType() string {
	return "Milestone"
}
