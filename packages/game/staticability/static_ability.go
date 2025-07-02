package staticability

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type StaticAbility interface {
	Name() string
	StaticKeyword() mtg.StaticKeyword
}

func New(staticAbilityDefinition *definition.StaticAbility) (StaticAbility, error) {
	cost, err := cost.Parse(staticAbilityDefinition.Cost)
	if err != nil {
		return nil, err
	}
	switch staticAbilityDefinition.Name {
	case "Cipher":
		return &Cipher{}, nil
	case "Flashback":
		return &Flashback{Cost: cost}, nil
	case "Replicate":
		return &Replicate{Cost: cost}, nil
	case "Splice":
		return NewSplice(cost, staticAbilityDefinition.Modifiers)
	default:
		return nil, fmt.Errorf("unknown static ability: %s", staticAbilityDefinition.Name)
	}
}
