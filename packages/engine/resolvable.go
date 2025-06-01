package engine

import "deckronomicon/packages/game/player"

// TODO this also lives in stack.go
type Resolvable interface {
	Description() string
	ID() string
	Name() string
	Resolve(*GameState, *player.Player) error
}
