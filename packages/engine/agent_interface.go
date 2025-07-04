package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/state"
)

// Need to pass in "Known game state to player" instead of "Game state".
// This is because the player would not know the full game state,
type PlayerAgent interface {
	GetNextAction(*state.Game) (Action, error)
	ReportState(*state.Game)
	Choose(*state.Game, choose.ChoicePrompt) (choose.ChoiceResults, error)
	PlayerID() string
}
