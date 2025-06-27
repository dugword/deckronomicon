package mtg

import (
	"errors"
	"fmt"
)

// Errors for game logic and actions.

var ErrAlreadyTapped = errors.New("object already tapped")

// Sentinel error
var ErrLibraryEmpty = errors.New("library empty")
var ErrGameOver = errors.New("game over")

var ErrInvalidZone = errors.New("invalid zone")

var ErrLandAlreadyPlayed = errors.New("land has already been played this turn")

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
