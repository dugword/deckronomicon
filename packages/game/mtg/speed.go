package mtg

type Speed string

const (
	SpeedInstant Speed = "Instant"
	SpeedSorcery Speed = "Sorcery"
)

func StringToSpeed(s string) (Speed, bool) {
	switch s {
	case "Instant":
		return SpeedInstant, true
	case "Sorcery":
		return SpeedSorcery, true
	default:
		return "", false
	}
}
