package event

func NewPassPriorityEvent(playerID string) GameEvent {
	return GameEvent{
		Type:     "PassPriority",
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
