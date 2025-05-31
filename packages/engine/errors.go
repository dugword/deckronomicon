package engine

import (
	"deckronomicon/packages/game/mtg"
	"errors"
	"fmt"
)

var ErrMaxTurnsExceeded = fmt.Errorf("maximum turns exceeded: %w", mtg.ErrGameOver)

// Standard errors
var ErrInvalidObjectType = errors.New("invalid object type")
var ErrObjectNotFound = errors.New("object not found")

// TODO: Define standard errors here for things like "no actions available"
// and also define custom error types for things like InvalidAction so I know
// if I can recover from them or not

// Game errors
var ErrInvalidAction = errors.New("invalid action")
