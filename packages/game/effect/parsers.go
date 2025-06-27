package effect

import (
	"deckronomicon/packages/game/mtg"
	"fmt"
)

func parseCount(modifiers map[string]any) (int, error) {
	count, ok := modifiers["Count"].(int)
	if !ok || count <= 0 {
		return 0, fmt.Errorf("a 'Count' modifier of type int greater than 0, got %T", modifiers["Count"])
	}
	return count, nil
}

func parseTargetPermanent(modifiers map[string]any) (mtg.TargetType, error) {
	targetString, ok := modifiers["Target"].(string)
	if ok && targetString != "Permanent" {
		return mtg.TargetTypeNone, fmt.Errorf("a 'Target' modifier of either empty or 'Permanent' required, got %q", targetString)
	}
	target, ok := mtg.StringToTargetType(targetString)
	if !ok {
		return mtg.TargetTypeNone, fmt.Errorf("invalid target type %q for Target modifier", targetString)
	}
	return target, nil
}

func parseTargetPlayer(modifiers map[string]any) (mtg.TargetType, error) {
	targetString, ok := modifiers["Target"].(string)
	if ok && targetString != "Player" {
		return mtg.TargetTypeNone, fmt.Errorf("a 'Target' modifier of either empty or 'Player' required, got %q", targetString)
	}
	if targetString == "" {
		return mtg.TargetTypeNone, nil
	}
	target, ok := mtg.StringToTargetType(targetString)
	if !ok {
		return mtg.TargetTypeNone, fmt.Errorf("invalid target type %q for Target modifier", targetString)
	}
	return target, nil
}

type query struct {
	CardTypes  []mtg.CardType `json:"CardTypes,omitempty"`
	Colors     []mtg.Color    `json:"Colors,omitempty"`
	Subtypes   []mtg.Subtype  `json:"Subtypes,omitempty"`
	ManaValues []int          `json:"ManaValues,omitempty"`
}

func parseQuery(modifiers map[string]any) (query, error) {
	var cardTypes []mtg.CardType
	cardTypesRaw, ok := modifiers["CardTypes"].([]any)
	if ok {
		for _, cardTypeRaw := range cardTypesRaw {
			cardType, ok := mtg.StringToCardType(fmt.Sprintf("%v", cardTypeRaw))
			if !ok {
				return query{}, fmt.Errorf("a 'CardTypes' modifier of type []CardType required, got %T", cardTypeRaw)
			}
			cardTypes = append(cardTypes, cardType)
		}
	}
	var colors []mtg.Color
	colorsRaw, ok := modifiers["Colors"].([]any)
	if ok {
		for _, colorRaw := range colorsRaw {
			color, ok := mtg.StringToColor(fmt.Sprintf("%v", colorRaw))
			if !ok {
				return query{}, fmt.Errorf("a 'Colors' modifier of type []Color required, got %T", colorRaw)
			}
			colors = append(colors, color)
		}
	}
	var subtypes []mtg.Subtype
	subtypesRaw, ok := modifiers["Subtypes"].([]any)
	if ok {
		for _, subtypeRaw := range subtypesRaw {
			subtype, ok := mtg.StringToSubtype(fmt.Sprintf("%v", subtypeRaw))
			if !ok {
				return query{}, fmt.Errorf("a 'Subtypes' modifier of type []Subtype required, got %T", subtypeRaw)
			}
			subtypes = append(subtypes, subtype)
		}
	}
	var manaValues []int
	manaValuesRaw, ok := modifiers["ManaValues"].([]any)
	if ok {
		for _, manaValueRaw := range manaValuesRaw {
			manaValue, ok := manaValueRaw.(int)
			if !ok {
				return query{}, fmt.Errorf("a 'ManaValues' modifier of type []int required, got %T", manaValueRaw)
			}
			manaValues = append(manaValues, manaValue)
		}
	}
	return query{
		CardTypes:  cardTypes,
		Colors:     colors,
		Subtypes:   subtypes,
		ManaValues: manaValues,
	}, nil
}
