package effect

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"fmt"
)

type LookAndChoose struct {
	CardTypes []mtg.CardType
	Choose    int
	Look      int
	Order     string
	Rest      mtg.Zone
}

func (e *LookAndChoose) Name() string {
	return "LookAndChoose"
}

func NewLookAndChoose(modifiers map[string]any) (*LookAndChoose, error) {
	look, ok := modifiers["Look"].(int)
	if !ok || look <= 0 {
		return nil, fmt.Errorf("a 'Look' modifier of type int greater than 0 required, got %T", modifiers["Look"])
	}
	choose, ok := modifiers["Choose"].(int)
	if !ok || choose <= 0 {
		return nil, fmt.Errorf("a 'Choose' modifier of type int greater than 0 required, got %T", modifiers["Choose"])
	}
	cardTypeVals, ok := modifiers["CardTypes"].([]any)
	if !ok || len(cardTypeVals) == 0 {
		return nil, fmt.Errorf("a non-empty 'CardTypes' modifier of type []mtg.CardType required, got %T", modifiers["CardTypes"])
	}
	var cardTypes []mtg.CardType
	for _, cardTypeVal := range cardTypeVals {
		cardTypeString, ok := cardTypeVal.(string)
		if !ok {
			return nil, fmt.Errorf("a 'CardTypes' modifier value required to be of type string, got %T", cardTypeVal)
		}
		cardType, ok := mtg.StringToCardType(cardTypeString)
		if !ok {
			return nil, fmt.Errorf("a 'CardTypes' modifier with CardType values required, got %q", cardTypeString)
		}
		cardTypes = append(cardTypes, cardType)
	}
	restString, ok := modifiers["Rest"].(string)
	if !ok {
		return nil, fmt.Errorf("a 'Rest' modifier of type string required, got %T", modifiers["Rest"])
	}
	rest, ok := mtg.StringToZone(restString)
	if !ok && (rest != mtg.ZoneLibrary && rest != mtg.ZoneGraveyard) {
		return nil, fmt.Errorf("a 'Rest' modifier of type Zone with value 'Library' or 'Graveyard' required, got %T", modifiers["Rest"])
	}
	order, ok := modifiers["Order"].(string)
	if ok && rest == mtg.ZoneLibrary && order != "Any" && order != "Random" {
		return nil, fmt.Errorf("a 'Order' modifier of type string with value 'Any' or 'Random' when rest is Library required, got %T", modifiers["Order"])

	}
	if look < choose {
		return nil, fmt.Errorf("a 'Look' (%d) value required to be greater than or equal to 'Choose' (%d)", look, choose)
	}
	return &LookAndChoose{
		CardTypes: cardTypes,
		Choose:    choose,
		Look:      look,
		Order:     order,
		Rest:      rest,
	}, nil
}

func (e *LookAndChoose) TargetSpec() target.TargetSpec {
	return target.NoneTargetSpec{}
}
