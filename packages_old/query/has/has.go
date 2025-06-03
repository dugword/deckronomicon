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
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			for _, cardType := range cardObj.CardTypes() {
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
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		objColors := cardObj.Colors()
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

func ID(id string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.ID() == id
	}
}

func Name(name string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.Name() == name
	}
}

func StaticKeyword(keywords ...mtg.StaticKeyword) query.Predicate {
	return func(obj query.Object) bool {
		for _, keyword := range keywords {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			for _, objKeyword := range cardObj.StaticKeywords() {
				if objKeyword == keyword {
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

func Subtype(subtypes ...mtg.Subtype) query.Predicate {
	return func(obj query.Object) bool {
		for _, subtype := range subtypes {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			for _, objSubtype := range cardObj.Subtypes() {
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
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			for _, objSupertype := range cardObj.Supertypes() {
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
