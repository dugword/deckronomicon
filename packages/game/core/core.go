package core

import "deckronomicon/packages/query"

type State interface {
	CanCastSorcery(string) bool
	GetNextID() string
	// Players() []Player
}

type Agent interface {
	ReportState(State)
}

type Ability interface {
	ID() string
	Name() string
	Resolve(state State, player Player) error
}

// TODO: Maybe move these into the resolve/spell/effect packages and split out
// the interfaces into smaller ones like we did for add to battlefield
type Player interface {
	// Agent() Agent
	AddMana(string) error
	DiscardCard(string) error
	Hand() query.View
	ID() string
	Life() int
	LoseLife(int) error
}

type Card interface {
	ID() string
	Name() string
	Match(query.Predicate) bool
}

type Object interface {
	ID() string
	Name() string
}

type Permanent interface {
	IsTapped() bool
	Tap() error
}

type Effect interface {
	ID() string
	Description() string
	Tags() []Tag
	// RequiresChoices(State, Player) bool
	// GetChoices(State, Player) ([]choose.ChoicePrompt, error)
	// Apply(State, Player, []choose.ChoiceResponse) error
	Apply(State, Player) error
}

type Tag struct {
	Key   string
	Value string
}
