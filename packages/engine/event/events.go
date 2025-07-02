package event

// TODO: Maybe have 2 packages, one for engine events and one for game events.
// Engine events are used to manage the game state, while game events are used
// for effect triggers and game rules. E.g. LandTappedForManaEvent would be a
// game event that is used to trigger effects that care about lands being tapped
// for mana, but AddMana would be a engine event that is used to add mana to
// the player's mana pool.

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

func (e *NoOpEvent) EventType() string {
	return "NoOp"
}

// MilestoneEvent is used to mark significant points in the game as defined by the player,
// this is used for analytics or tracking purposes.
type MilestoneEvent struct {
	Message string
}

func (e *MilestoneEvent) EventType() string {
	return "Milestone"
}
