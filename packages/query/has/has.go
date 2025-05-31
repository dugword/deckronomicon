package has

import (
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"fmt"
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

func CardType(cardTypes ...mtg.CardType) query.Predicate {
	return func(obj query.Object) bool {
		for _, t := range cardTypes {
			fmt.Println("Checking card type", t)
			fmt.Printf("Object is %T\n", obj)
			found := false
			cardObj, ok := obj.(query.CardObject)
			if !ok {
				fmt.Println("Card was not a card object")
				return false
			}
			fmt.Println("Card was a card object")
			if slices.Contains(cardObj.CardTypes(), t) {
				found = true
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

func Controller(controllerID string) query.Predicate {
	return func(obj query.Object) bool {
		return obj.Controller() == controllerID
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

func ManaValue(value int) query.Predicate {
	return func(obj query.Object) bool {
		cardObj, ok := obj.(query.CardObject)
		if !ok {
			return false
		}
		return cardObj.ManaValue() == value
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

func Subtype(subtypes ...mtg.Subtype) query.Predicate {
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

func Supertype(supertypes ...mtg.Supertype) query.Predicate {
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
