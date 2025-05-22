package game

import (
	"errors"
	"fmt"
)

// Sentinel error
var ErrLibraryEmpty = errors.New("library empty")
var ErrGameOver = errors.New("game over")
var ErrMaxTurnsExceeded = fmt.Errorf("maximum turns exceeded: %w", ErrGameOver)

// Standard errors
var ErrInvalidObjectType = errors.New("invalid object type")
var ErrObjectNotFound = errors.New("object not found")

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
