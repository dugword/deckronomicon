package is

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

func Land() query.Predicate {
	return func(obj query.Object) bool {
		for _, cardType := range obj.CardTypes() {
			if cardType == mtg.CardTypeLand {
				return true
			}
		}
		return false
	}
}

func Not(predicate query.Predicate) query.Predicate {
	return func(obj query.Object) bool {
		return !predicate(obj)
	}
}

func Permanent() query.Predicate {
	return func(obj query.Object) bool {
		for _, cardType := range obj.CardTypes() {
			if cardType.IsPermanent() {
				return true
			}
		}
		return false
	}
}
