package event

// Events go in the event record log so they can be replayed later. That means
// the properties should be serializable and public.

const (
	EventTypeDrawCard = "DrawCard"
)

// TODO: maybe use typed constants for event types
const (
	EventTypeSetNextPlayer = "SetNextPlayer"
)

type GameEvent interface {
	EventType() string
}

type Source interface {
	Name() string
}

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
