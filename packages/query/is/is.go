package is

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
)

func Land() query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		for _, cardType := range cardObj.CardTypes() {
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
		fmt.Println("card is not object")
		for _, cardType := range cardObj.CardTypes() {
			if cardType.IsPermanent() {
				return true
			}
		}
		return false
	}
}
