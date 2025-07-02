package is

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"slices"
)

func Spell() query.Predicate {
	return func(obj query.Object) bool {
		_, ok := obj.(*gob.Spell)
		return ok
	}
}

func AbilityOnStack() query.Predicate {
	return func(obj query.Object) bool {
		_, ok := obj.(*gob.AbilityOnStack)
		return ok
	}
}

func Card() query.Predicate {
	return func(obj query.Object) bool {
		_, ok := obj.(query.CardObject)
		return ok
	}
}

func Not(predicate query.Predicate) query.Predicate {
	return func(obj query.Object) bool {
		return !predicate(obj)
	}
}

func Permanent() query.Predicate {
	return func(obj query.Object) bool {
		_, ok := obj.(*gob.Permanent)
		return ok
	}
}

func Tapped() query.Predicate {
	return func(obj query.Object) bool {
		permObj, ok := obj.(*gob.Permanent)
		if !ok {
			return false
		}
		if !permObj.IsTapped() {
			return false
		}
		return true
	}
}

func PermanentCardType() query.Predicate {
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

func Land() query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.CardTypes(), mtg.CardTypeLand)
	}
}

func Untapped() query.Predicate {
	return func(obj query.Object) bool {
		permObj, ok := obj.(*gob.Permanent)
		if !ok {
			return false
		}
		return !permObj.IsTapped()
	}
}

func ManaAbility() query.Predicate {
	return func(obj query.Object) bool {
		if abilityObj, ok := obj.(*gob.Ability); ok {
			for _, efct := range abilityObj.Effects() {
				if _, ok := efct.(*effect.AddMana); ok {
					return true
				}
			}
		}
		return false
	}
}
