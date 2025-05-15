package game

import (
	"errors"
	"fmt"
)

// Sentinel error
var ErrGameOver = errors.New("game over")

var ErrMaxTurnsExceeded = fmt.Errorf("maximum turns exceeded: %w", ErrGameOver)

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

func (e PlayerLostError) Error() string {
	return fmt.Sprintf("player lost: %s", e.Reason)
}

func (e PlayerLostError) Unwrap() error {
	return ErrGameOver
}

// DescribeGameOver returns a friendly message about why the game ended
/*
func DescribeGameOver(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, ErrGameOver) {
		var ple PlayerLostError
		if errors.As(err, &ple) {
			return fmt.Sprintf("Player %s lost the game (%s)", ple.PlayerID, ple.Reason.String())
		}

	}

	return ""
}
*/
