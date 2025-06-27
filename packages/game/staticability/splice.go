package staticability

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Splice struct {
	Cost    cost.Cost
	Subtype mtg.Subtype `json:"Subtype"`
}

func (a Splice) Name() string {
	return "Splice"
}

func (a Splice) StaticKeyword() mtg.StaticKeyword {
	return mtg.StaticKeywordSplice
}

func NewSplice(cost cost.Cost, modifiers map[string]any) (Splice, error) {
	subtypeString, ok := modifiers["Subtype"].(string)
	if !ok {
		return Splice{}, fmt.Errorf("a 'Subtype' key is required for Splice modifier")
	}
	subtype, ok := mtg.StringToSubtype(subtypeString)
	if !ok {
		return Splice{}, fmt.Errorf("invalid subtype %s for Splice modifier", subtypeString)
	}
	return Splice{
		Cost:    cost,
		Subtype: subtype,
	}, nil
}
