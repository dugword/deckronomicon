package engine

import "deckronomicon/packages/choose"

type TurnBasedAction interface {
	Description() string
	RequiresChoice() bool
	GetChoices(state *GameState) ([]choose.Choice, error)
	Resolve(state *GameState, responses []choose.ChoiceResponse) (string, error)
}
