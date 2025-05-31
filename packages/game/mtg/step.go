package mtg

import "fmt"

type Step string

// TODO: Standardize on all the types and string to type and use a standard
// error

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

// StringToCardType converts a string to a CardType.
func StringToStep(s string) (Step, error) {
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
		return "", fmt.Errorf("unknown step: %s", s)
	}
	return cardType, nil
}
