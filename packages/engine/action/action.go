package action

import (
	"deckronomicon/packages/engine"
)

type Request interface {
	Build(playerID string) (engine.Action, error)
	Name() string
}
