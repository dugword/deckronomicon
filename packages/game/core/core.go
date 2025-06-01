package core

import "deckronomicon/packages/query"

type State interface {
	GetNextID() string
	CanCastSorcery(string) bool
	// Players() []Player
}

type Agent interface {
	ReportState(State)
}

type Player interface {
	// Agent() Agent
	AddMana(string) error
	DiscardCard(string) error
	Hand() query.View
	Life() int
	LoseLife(int) error
}

type Card interface {
	ID() string
	Name() string
}

type Object interface {
	Name() string
}

type Permanent interface {
	IsTapped() bool
	Tap() error
}
