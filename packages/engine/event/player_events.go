package event

import "deckronomicon/packages/mana"

const (
	EventTypeAddMana = "AddMana"
)

type PlayerEvent interface {
	isPlayerEvent()
}

type PlayerBaseEvent struct {
	PlayerID string
}

func (e PlayerBaseEvent) isPlayerEvent() {}

type AddManaEvent struct {
	PlayerBaseEvent
	ManaType mana.ManaType
	Amount   int
}

func (e AddManaEvent) EventType() string {
	return EventTypeAddMana
}

func NewAddManaEvent(playerID string, manaType mana.ManaType, amount int) AddManaEvent {
	return AddManaEvent{
		PlayerBaseEvent: PlayerBaseEvent{
			PlayerID: playerID,
		},
		ManaType: manaType,
		Amount:   amount,
	}
}

/*
func NewPassPriorityEvent(playerID string) GameEvent {
	return GameEvent{
		PlayerID: playerID,
	}
}

func NewConcedeGameEvent(playerID string) GameEvent {
	return GameEvent{}
}

func NewChooseInitialHandEvent(playerID string, cardIDs []string) GameEvent {
	return GameEvent{}
}

func NewMulliganEvent(playerID string, cardsMulliganed []string) GameEvent {
	return GameEvent{}
}
*/
