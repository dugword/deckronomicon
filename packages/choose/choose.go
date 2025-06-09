package choose

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

type Source interface {
	Name() string
}

type Choice interface {
	Name() string
	ID() string
	Zone() mtg.Zone
}

// TODO: Use an interface for ChoiceResult
type ChoicePrompt2 interface {
	PromptID() string
	Type() string
	Choices(game state.Game) ([]Choice, error)
	Validate(input string) error
	Apply(choice string, game state.Game) (ResultStruct, error)
}

/*
type ChooseCardPrompt struct {
	Zone     string
	PlayerID string
	Count    int
}
*/

type ChooseTargetPrompt struct {
}

type ResultStruct struct {
	SelectedIDs []string
}

type ChoicePrompt struct {
	Choices    []Choice
	MaxChoices int
	Message    string
	MinChoices int
	Optional   bool
	Source     Source
}
