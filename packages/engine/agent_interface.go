package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/state"
)

type PlayerAgent interface {
	GetNextAction(*state.Game) (Action, error)
	ReportState(*state.Game)
	Choose(choose.ChoicePrompt) (choose.ChoiceResults, error)
	PlayerID() string
}
