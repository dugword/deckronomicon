package mtg

type Duration string

const (
	DurationEndOfTurn Duration = "EndOfTurn"
)

func StringToDuration(s string) (Duration, bool) {
	switch s {
	case string(DurationEndOfTurn):
		return DurationEndOfTurn, true
	default:
		return "", false
	}
}
