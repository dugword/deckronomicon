package effect

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"fmt"
	"strconv"
)

func buildQuery(
	modifiers []definition.EffectModifier,
) (query.Predicate, error) {
	var cardTypes []mtg.CardType
	var colors []mtg.Color
	var subtypes []mtg.Subtype
	var manaValues []int
	for _, modifier := range modifiers {
		if modifier.Key == "CardType" {
			cardType, ok := mtg.StringToCardType(modifier.Value)
			if !ok {
				return nil, fmt.Errorf("invalid card type %q for Search effect", modifier.Value)
			}
			fmt.Println("Adding card type to query", cardType)
			cardTypes = append(cardTypes, cardType)
		}
		if modifier.Key == "Color" {
			color, ok := mtg.StringToColor(modifier.Value)
			if !ok {
				return nil, fmt.Errorf("invalid color %q for Search effect", modifier.Value)
			}
			colors = append(colors, color)
		}
		if modifier.Key == "Subtype" {
			subtype, ok := mtg.StringToSubtype(modifier.Value)
			if !ok {
				return nil, fmt.Errorf("invalid subtype %q for Search effect", modifier.Value)
			}
			subtypes = append(subtypes, subtype)
		}
		if modifier.Key == "ManaValue" {
			manaValue, err := strconv.Atoi(modifier.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid mana value %q for Search effect: %w", modifier.Value, err)
			}
			manaValues = append(manaValues, manaValue)
		}
	}
	// TODO: Can this be more simple/elegant?
	var cardTypePredicates []query.Predicate
	for _, cardType := range cardTypes {
		cardTypePredicates = append(cardTypePredicates, has.CardType(cardType))
	}
	var colorPredicates []query.Predicate
	for _, color := range colors {
		colorPredicates = append(colorPredicates, has.Color(color))
	}
	var subtypePredicates []query.Predicate
	for _, subtype := range subtypes {
		subtypePredicates = append(subtypePredicates, has.Subtype(subtype))
	}
	var manaValuePredicates []query.Predicate
	for _, manaValue := range manaValues {
		manaValuePredicates = append(manaValuePredicates, has.ManaValue(manaValue))
	}
	var predicates []query.Predicate
	if len(cardTypePredicates) != 0 {
		predicates = append(predicates, query.Or(cardTypePredicates...))
	}
	if len(colorPredicates) != 0 {
		predicates = append(predicates, query.Or(colorPredicates...))
	}
	if len(subtypePredicates) != 0 {
		predicates = append(predicates, query.Or(subtypePredicates...))
	}
	if len(manaValuePredicates) != 0 {
		predicates = append(predicates, query.Or(manaValuePredicates...))
	}
	query := query.And(predicates...)
	return query, nil
}
