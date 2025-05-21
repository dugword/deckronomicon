package game

// ChoiceResolver is an interface for resolving player choices.
// Maybe do something where I can pass in "play Island" and it'll take the second param as the Choice and only prompt if it is missing
// maybe support typing in the number or the name of the card
type ChoiceResolver interface {
	// TODO: Source needs formating or a struct or something...
	// could be an action, or a spell, or a card name, or an ability
	ChooseOne(prompt string, source string, choices []Choice) (Choice, error)
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	//ChooseAny(prompt string, choices []Choice) []Choice
	//Confirm(prompt string) bool // For simple yes/no prompts
}

// PlayerAgent defines how player decisions are made.
type PlayerAgent interface {
	ChoiceResolver
	GetNextAction(state *GameState) GameAction
	ReportState(state *GameState)
}

// Choice represents a choice made by the player.
type Choice struct {
	Index  int
	Name   string
	Source string
}
