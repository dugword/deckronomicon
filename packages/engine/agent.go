package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/state"
)

type PlayerAgent interface {
	// TODO Will be a complex type in the future, string works for now
	GetNextAction(state.Game) (Action, error)
	ReportState(state.Game) error
	Choose(choose.ChoicePrompt) ([]choose.Choice, error)
}
