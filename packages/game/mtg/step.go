package mtg

type Step string

const (
	StepUntap             Step = "Untap"
	StepUpkeep            Step = "Upkeep"
	StepDraw              Step = "Draw"
	StepPrecombatMain     Step = "PrecombatMain"
	StepBeginningOfCombat Step = "BeginningOfCombat"
	StepDeclareAttackers  Step = "DeclareAttackers"
	StepDeclareBlockers   Step = "DeclareBlockers"
	StepCombatDamage      Step = "CombatDamage"
	StepEndOfCombat       Step = "EndOfCombat"
	StepPostcombatMain    Step = "PostcombatMain"
	StepEnd               Step = "End"
	StepCleanup           Step = "Cleanup"
)

func StringToStep(s string) (Step, bool) {
	stringToStep := map[string]Step{
		"Untap":             StepUntap,
		"Upkeep":            StepUpkeep,
		"Draw":              StepDraw,
		"PrecombatMain":     StepPrecombatMain,
		"BeginningOfCombat": StepBeginningOfCombat,
		"DeclareAttackers":  StepDeclareAttackers,
		"DeclareBlockers":   StepDeclareBlockers,
		"CombatDamage":      StepCombatDamage,
		"EndOfCombat":       StepEndOfCombat,
		"PostcombatMain":    StepPostcombatMain,
		"End":               StepEnd,
		"Cleanup":           StepCleanup,
	}
	cardType, ok := stringToStep[s]
	if !ok {
		return "", false
	}
	return cardType, true
}
