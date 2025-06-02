package core

import (
	"deckronomicon/packages/query"
)

type State interface {
	CanCastSorcery(string) bool
	GetNextID() string
	// Players() []Player
}

type Agent interface {
	ReportState(State)
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
}

type Object interface {
	ID() string
	Name() string
}

type Permanent interface {
	IsTapped() bool
	Tap() error
}

type Tag struct {
	Key   string
	Value string
}
