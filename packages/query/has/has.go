package has

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

func Any(predicates ...query.Predicate) query.Predicate {
	return func(obj query.Object) bool {
		for _, p := range predicates {
			if p(obj) {
				return true
			}
		}
		return false
	}
}

func All(predicates ...query.Predicate) query.Predicate {
	return func(obj query.Object) bool {
		for _, p := range predicates {
			if !p(obj) {
				return false
			}
		}
		return true
	}
}

func CardType(cardTypes ...mtg.CardType) query.Predicate {
	return func(obj query.Object) bool {
		for _, t := range cardTypes {
			found := false
			for _, cardType := range obj.CardTypes() {
				if cardType == t {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

// TODO support multiple colors
func Color(color mtg.Color) query.Predicate {
	return func(obj query.Object) bool {
		objColors := obj.Colors()
		switch color {
		case mtg.ColorBlack:
			return objColors.Black
		case mtg.ColorBlue:
			return objColors.Blue
		case mtg.ColorGreen:
			return objColors.Green
		case mtg.ColorRed:
			return objColors.Red
		case mtg.ColorWhite:
			return objColors.White
		default:
			return false
		}
	}
}

func Name(name string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.Name() == name
	}
}

/*
func StaticAbilities(abilities ...mtg.StaticKeyword) predicate.Predicate {
	return func(obj query.Object) bool {
		for _, ability := range abilities {
			found := false
			for _, objAbility := range obj.StaticAbilities() {
				if objAbility == ability {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}
*/

func Subtype(subtypes ...mtg.Subtype) query.Predicate {
	return func(obj query.Object) bool {
		for _, subtype := range subtypes {
			found := false
			for _, objSubtype := range obj.Subtypes() {
				if objSubtype == subtype {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func Supertype(supertypes ...mtg.Supertype) query.Predicate {
	return func(obj query.Object) bool {
		for _, supertype := range supertypes {
			found := false
			for _, objSupertype := range obj.Supertypes() {
				if objSupertype == supertype {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}
