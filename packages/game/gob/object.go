package gob

import "deckronomicon/packages/query"

type Object interface {
	Controller() string
	Description() string
	ID() string
	Name() string
	Owner() string
	Match(query.Predicate) bool
}
