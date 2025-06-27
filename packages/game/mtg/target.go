package mtg

type TargetType string

const (
	TargetTypeNone      TargetType = "None"
	TargetTypePlayer    TargetType = "Player"
	TargetTypeCreature  TargetType = "Creature"
	TargetTypeSpell     TargetType = "Spell"
	TargetTypePermanent TargetType = "Permanent"
)

func StringToTargetType(s string) (TargetType, bool) {
	targetTypes := map[string]TargetType{
		string(TargetTypeNone):      TargetTypeNone,
		string(TargetTypePlayer):    TargetTypePlayer,
		string(TargetTypeCreature):  TargetTypeCreature,
		string(TargetTypeSpell):     TargetTypeSpell,
		string(TargetTypePermanent): TargetTypePermanent,
	}
	val, ok := targetTypes[s]
	return val, ok
}
