package mtg

type Phase string

const (
	PhaseBeginning      Phase = "Beginning"
	PhasePrecombatMain  Phase = "PrecombatMain"
	PhaseCombat         Phase = "Combat"
	PhasePostcombatMain Phase = "PostcombatMain"
	PhaseEnding         Phase = "Ending"
)

func StringToPhase(s string) (Phase, bool) {
	stringToPhase := map[string]Phase{
		"Beginning":      PhaseBeginning,
		"PrecombatMain":  PhasePrecombatMain,
		"Combat":         PhaseCombat,
		"PostcombatMain": PhasePostcombatMain,
		"Ending":         PhaseEnding,
	}
	phase, ok := stringToPhase[s]
	if !ok {
		return "", false
	}
	return phase, true
}
