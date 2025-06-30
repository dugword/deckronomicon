package mtg

type TargetType string

const (
	TargetTypeNone      TargetType = "None"
	TargetTypePlayer    TargetType = "Player"
	TargetTypeCard      TargetType = "Card"
	TargetTypeSpell     TargetType = "Spell"
	TargetTypePermanent TargetType = "Permanent"
)

// TODO: This is a good pattern and I think I need to change how
// I manage card types/subtypes/colors to use this pattern as well.
func StringToTargetType(s string) (TargetType, bool) {
	targetTypes := map[string]TargetType{
		string(TargetTypeNone):      TargetTypeNone,
		string(TargetTypePlayer):    TargetTypePlayer,
		string(TargetTypeCard):      TargetTypeCard,
		string(TargetTypeSpell):     TargetTypeSpell,
		string(TargetTypePermanent): TargetTypePermanent,
	}
	val, ok := targetTypes[s]
	return val, ok
}
