package event

const (
	EventTypeTapPermanent = "TapPermanent"
)

type PermanentEvent interface {
	isPermanentEvent()
}

type PermanentBaseEvent struct {
	PermanentID string
}

func (e PermanentBaseEvent) isPermanentEvent() {}

type TapPermanentEvent struct {
	PermanentBaseEvent
}

func (e TapPermanentEvent) EventType() string {
	return EventTypeTapPermanent
}
func (e TapPermanentEvent) isPermanentEvent() {}

func NewTapPermanentEvent(permanentID string) TapPermanentEvent {
	return TapPermanentEvent{
		PermanentBaseEvent: PermanentBaseEvent{
			PermanentID: permanentID,
		},
	}
}

/*

func NewUntapPermanentEvent(playerID, permID string) GameEvent {
	return GameEvent{}
}

func NewPermanentDiesEvent(permID string, reason string) GameEvent {
	return GameEvent{}
}

func NewMoveToZoneEvent(cardID string, from, to string) GameEvent {
	return GameEvent{}
}
*/
