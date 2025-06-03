package mtg

type Phase string

const (
	PhaseBeginning      Phase = "Beginning"
	PhasePrecombatMain  Phase = "PrecombatMain"
	PhaseCombat         Phase = "Combat"
	PhasePostcombatMain Phase = "PostcombatMain"
	PhaseEnding         Phase = "Ending"
)
