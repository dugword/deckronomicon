package effect

import (
	"deckronomicon/packages/game/target"
	"fmt"
)

type AddMana struct {
	Mana string
}

func (e AddMana) Name() string {
	return "AddMana"
}

func NewAddMana(modifiers map[string]any) (AddMana, error) {
	manaString, ok := modifiers["Mana"].(string)
	if !ok {
		return AddMana{}, fmt.Errorf("invalid 'Mana' modifier %q", modifiers)
	}
	return AddMana{Mana: manaString}, nil
}

// Move this to not the resolver, but the effect description like with the modifiers
func (e AddMana) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
