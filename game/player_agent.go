package game

// PlayerAgent defines how player decisions are made.
type PlayerAgent interface {
	// TODO: Not sure I love this here
	//ChooseAny(prompt string, choices []Choice) []Choice
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	ChooseMany(prompt string, source ChoiceSource, choices []Choice) ([]Choice, error)
	ChooseOne(prompt string, source ChoiceSource, choices []Choice) (Choice, error)
	Confirm(prompt string, source ChoiceSource) (bool, error)
	EnterNumber(prompt string, source ChoiceSource) (int, error)
	GetNextAction(state *GameState) *GameAction
	PlayerID() string
	ReportState(state *GameState)
}
