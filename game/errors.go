package game

import (
	"errors"
	"fmt"
)

// TODO: Define standard errors here for things like "no actions available"
// and also define custom error types for things like InvalidAction so I know
// if I can recover from them or not

// Game errors
var ErrInvalidAction = errors.New("invalid action")

// Sentinel error
var ErrLibraryEmpty = errors.New("library empty")
var ErrGameOver = errors.New("game over")
var ErrMaxTurnsExceeded = fmt.Errorf("maximum turns exceeded: %w", ErrGameOver)

// Standard errors
var ErrInvalidObjectType = errors.New("invalid object type")
var ErrObjectNotFound = errors.New("object not found")

// TODO I don't like this name, come up with something better
var ErrAlreadyTapped = errors.New("object already tapped")

// PlayerLostReason provides detailed cause for game loss
type PlayerLostReason string

const (
	LifeTotalZero PlayerLostReason = "life total zero"
	DeckedOut     PlayerLostReason = "decked out"
	Conceded      PlayerLostReason = "player conceded"
)

// PlayerLostError indicates a player has lost
type PlayerLostError struct {
	Reason PlayerLostReason
}

// Error returns the error message for PlayerLostError
func (e PlayerLostError) Error() string {
	return fmt.Sprintf("player lost: %s", e.Reason)
}

// Unwrap returns the underlying error for PlayerLostError
func (e PlayerLostError) Unwrap() error {
	return ErrGameOver
}
