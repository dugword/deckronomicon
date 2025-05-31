package player

import (
	"deckronomicon/packages/choice"
	"deckronomicon/packages/game/action"
)

type GameState any

// PlayerAgent defines how player decisions are made.
type Agent interface {
	// TODO: Not sure I love this here
	//ChooseAny(prompt string, choices []Choice) []Choice
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	ChooseMany(prompt string, source choice.Source, choices []choice.Choice) ([]choice.Choice, error)
	ChooseOne(prompt string, source choice.Source, choices []choice.Choice) (choice.Choice, error)
	Confirm(prompt string, source choice.Source) (bool, error)
	EnterNumber(prompt string, source choice.Source) (int, error)
	GetNextAction(GameState) (action.Action, error)
	RegisterPlayer(*Player)
	ReportState(GameState)
}
