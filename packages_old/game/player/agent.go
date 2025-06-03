package player

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game"
)

type GameState any

// PlayerAgent defines how player decisions are made.
type Agent interface {
	// TODO: Not sure I love this here
	//ChooseAny(prompt string, choices []Choice) []Choice
	//ChooseN(prompt string, choices []Choice, n int) []Choice
	//ChooseUpToN(prompt string, choices []Choice, n int) []Choice
	ChooseMany(prompt string, source choose.Source, choices []choose.Choice) ([]choose.Choice, error)
	ChooseOne(prompt string, source choose.Source, choices []choose.Choice) (choose.Choice, error)
	Confirm(prompt string, source choose.Source) (bool, error)
	EnterNumber(prompt string, source choose.Source) (int, error)
	GetNextAction(GameState) (game.Action, error)
	RegisterPlayer(*Player)
	ReportState(GameState)
}
