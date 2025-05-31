package is

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"slices"
)

func Land() query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.CardTypes(), mtg.CardTypeLand)
	}
}

// Spell returns true if the card is a spell. Spells are cards that are not lands.
func Spell() query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return !slices.Contains(cardObj.CardTypes(), mtg.CardTypeLand)
	}
}

func Not(predicate query.Predicate) query.Predicate {
	return func(obj query.Object) bool {
		return !predicate(obj)
	}
}

func Tapped() query.Predicate {
	return func(obj query.Object) bool {
		permObj, ok := obj.(query.PermanentObject)
		if !ok {
			return false
		}
		if !permObj.IsTapped() {
			return false
		}
		return true
	}
}

func Permanent() query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		for _, cardType := range cardObj.CardTypes() {
			if cardType.IsPermanent() {
				return true
			}
		}
		return false
	}
}
