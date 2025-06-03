package event

type GameEvent struct {
	Action      string
	Choices     []string
	Phase       string
	PlayerID    string
	RNGSnapshot int64
	Result      string
	Source      string
	Targets     []string
	Turn        int
	Type        string
}

func NewGameOverEvent(winnerID string) GameEvent {
	return GameEvent{}
}
