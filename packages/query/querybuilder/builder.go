package querybuilder

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
)

// TODO: This is redundant with the one in effect package.
// Maybe we should have a common package for this?
// Or maybe this should be in the query package?

type Opts struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func Build(
	opts Opts,
) (query.Predicate, error) {
	// TODO: Can this be more simple/elegant?
	var cardTypePredicates []query.Predicate
	for _, cardType := range opts.CardTypes {
		cardTypePredicates = append(cardTypePredicates, has.CardType(cardType))
	}
	var colorPredicates []query.Predicate
	for _, color := range opts.Colors {
		colorPredicates = append(colorPredicates, has.Color(color))
	}
	var subtypePredicates []query.Predicate
	for _, subtype := range opts.Subtypes {
		subtypePredicates = append(subtypePredicates, has.Subtype(subtype))
	}
	var manaValuePredicates []query.Predicate
	for _, manaValue := range opts.ManaValues {
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
