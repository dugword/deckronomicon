package mtg

import "fmt"

type Speed string

const (
	SpeedInstant Speed = "Instant"
	SpeedSorcery Speed = "Sorcery"
)

func StringToSpeed(s string) (Speed, error) {
	switch s {
	case "Instant":
		return SpeedInstant, nil
	case "Sorcery":
		return SpeedSorcery, nil
	default:
		return "", fmt.Errorf("unknown Speed %q", s)
	}
}
