package engine

import "deckronomicon/packages/state"

type PlayerAgent interface {
	// TODO Will be a complex type in the future, string works for now
	GetNextAction() (Action, error)
	ReportState(state state.Game) error
}
