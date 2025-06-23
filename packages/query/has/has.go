package has

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"slices"
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

func AnyCardType(cardTypes ...mtg.CardType) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		for _, t := range cardTypes {
			if slices.Contains(cardObj.CardTypes(), t) {
				return true
			}
		}
		return false
	}
}

func CardType(cardType mtg.CardType) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.CardTypes(), cardType)
	}
}

func AllCardTypes(cardTypes ...mtg.CardType) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		for _, t := range cardTypes {
			if !slices.Contains(cardObj.CardTypes(), t) {
				return false
			}
		}
		return true
	}
}

func AnySourceID(sourceIDs ...string) query.Predicate {
	return func(obj query.Object) bool {
		stackObj, ok := obj.(query.StackObject)
		if !ok {
			return false
		}
		return slices.Contains(sourceIDs, stackObj.SourceID())
	}
}

func SourceID(sourceID string) query.Predicate {
	return func(obj query.Object) bool {
		stackObj, ok := obj.(query.StackObject)
		if !ok {
			return false
		}
		return stackObj.SourceID() == sourceID
	}
}

func AnyColors(colors ...mtg.Color) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		objColors := cardObj.Colors()
		for _, color := range colors {
			switch color {
			case mtg.ColorBlack:
				if objColors.Black {
					return true
				}
			case mtg.ColorBlue:
				if objColors.Blue {
					return true
				}
			case mtg.ColorGreen:
				if objColors.Green {
					return true
				}
			case mtg.ColorRed:
				if objColors.Red {
					return true
				}
			case mtg.ColorWhite:
				if objColors.White {
					return true
				}
			default:
				continue
			}
		}
		return false
	}
}

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

func AllColors(colors ...mtg.Color) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		objColors := cardObj.Colors()
		for _, color := range colors {
			switch color {
			case mtg.ColorBlack:
				if !objColors.Black {
					return false
				}
			case mtg.ColorBlue:
				if !objColors.Blue {
					return false
				}
			case mtg.ColorGreen:
				if !objColors.Green {
					return false
				}
			case mtg.ColorRed:
				if !objColors.Red {
					return false
				}
			case mtg.ColorWhite:
				if !objColors.White {
					return false
				}
			default:
				continue
			}
		}
		return true
	}
}

func AnyControllerID(controllerIDs ...string) query.Predicate {
	return func(obj query.Object) bool {
		return slices.Contains(controllerIDs, obj.Controller())
	}
}

func Controller(controllerID string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.Controller() == controllerID
	}
}

func AnyID(ids ...string) query.Predicate {
	return func(obj query.Object) bool {
		return slices.Contains(ids, obj.ID())
	}
}

func ID(id string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.ID() == id
	}
}

func AnyName(names ...string) query.Predicate {
	return func(obj query.Object) bool {
		return slices.Contains(names, obj.Name())
	}
}

func Name(name string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.Name() == name
	}
}

func AnyManaValue(values ...int) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(values, cardObj.ManaValue())
	}
}

func ManaValue(value int) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return cardObj.ManaValue() == value
	}
}

func AnyStaticKeyword(keywords ...mtg.StaticKeyword) query.Predicate {
	return func(obj query.Object) bool {
		for _, keyword := range keywords {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.StaticKeywords(), keyword) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func StaticKeyword(keyword mtg.StaticKeyword) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.StaticKeywords(), keyword)
	}
}

func AllStaticKeywords(keywords ...mtg.StaticKeyword) query.Predicate {
	return func(obj query.Object) bool {
		for _, keyword := range keywords {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.StaticKeywords(), keyword) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func AnySubtype(subtypes ...mtg.Subtype) query.Predicate {
	return func(obj query.Object) bool {
		for _, subtype := range subtypes {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.Subtypes(), subtype) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func Subtype(subtype mtg.Subtype) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.Subtypes(), subtype)
	}
}

func AllSubtypes(subtypes ...mtg.Subtype) query.Predicate {
	return func(obj query.Object) bool {
		for _, subtype := range subtypes {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.Subtypes(), subtype) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func AnySupertype(supertypes ...mtg.Supertype) query.Predicate {
	return func(obj query.Object) bool {
		for _, supertype := range supertypes {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.Supertypes(), supertype) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func Supertype(supertype mtg.Supertype) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return slices.Contains(cardObj.Supertypes(), supertype)
	}
}

func AllSupertypes(supertypes ...mtg.Supertype) query.Predicate {
	return func(obj query.Object) bool {
		for _, supertype := range supertypes {
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				return false
			}
			if slices.Contains(cardObj.Supertypes(), supertype) {
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
}
