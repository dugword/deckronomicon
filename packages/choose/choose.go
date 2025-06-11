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
	//Zone() mtg.Zone // TODO: Is Zone right? Maybe something more generic?
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

type GenericChoice struct {
	name string
	id   string
	zone mtg.Zone
}

func NewGenericChoice(name, id string, zone mtg.Zone) GenericChoice {
	return GenericChoice{
		name: name,
		id:   id,
		zone: zone,
	}
}

func (c GenericChoice) Name() string {
	return c.name
}

func (c GenericChoice) ID() string {
	return c.id
}

func (c GenericChoice) Zone() mtg.Zone {
	return c.zone
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
