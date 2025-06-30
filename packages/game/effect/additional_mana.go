package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type AdditionalMana struct {
	Subtype  mtg.Subtype
	Mana     string
	Duration mtg.Duration
}

func (e AdditionalMana) Name() string {
	return "AdditionalMana"
}

func NewAdditionalMana(modifiers map[string]any) (AdditionalMana, error) {
	subtypeString, ok := modifiers["Subtype"].(string)
	if !ok {
		return AdditionalMana{}, fmt.Errorf("a 'Subtype' modifier of type string required, got %T", modifiers["Subtype"])
	}
	subtype, ok := mtg.StringToSubtype(subtypeString)
	if !ok {
		return AdditionalMana{}, fmt.Errorf("a valid 'Subtype' modifier required, got %q", subtypeString)
	}
	manaString, ok := modifiers["Mana"].(string)
	if !ok {
		return AdditionalMana{}, fmt.Errorf("a 'Mana' modifier of type string required, got %T", modifiers["Mana"])
	}
	if manaString == "" {
		return AdditionalMana{}, fmt.Errorf("a non-empty 'Mana' modifier required")
	}
	durationString, ok := modifiers["Duration"].(string)
	if !ok {
		return AdditionalMana{}, fmt.Errorf("a 'Duration' modifier of type Duration required, got %T", modifiers["Duration"])
	}
	duration, ok := mtg.StringToDuration(durationString)
	if !ok {
		return AdditionalMana{}, fmt.Errorf("a valid 'Duration' modifier required, got %q", durationString)
	}
	if duration == "" {
		return AdditionalMana{}, fmt.Errorf("a non-empty 'Duration' modifier required")
	}
	return AdditionalMana{
		Subtype:  subtype,
		Mana:     manaString,
		Duration: duration,
	}, nil
}

func (e AdditionalMana) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
