package query

import (
	"deckronomicon/packages/game/mtg"
	"slices"
)

// TODO Maybe break this into more specific interfaces
type Object interface {
	Description() string
	ID() string
	Match(Predicate) bool
	Name() string
}

type CardObject interface {
	Object
	CardTypes() []mtg.CardType
	Colors() mtg.Colors
	StaticKeywords() []mtg.StaticKeyword
	Subtypes() []mtg.Subtype
	Supertypes() []mtg.Supertype
}

type AbilityObject interface {
}

type PermanentObject interface {
	CardObject
	IsTapped() bool
	HasSummoningSickness() bool
}

type Card interface{}

type Permanent interface{}

type Ability interface{}

// TODO add a ToString method to Predicate so we can have a string
// representation of the predicate for debugging purposes
type Predicate func(obj Object) bool

func And(predicates ...Predicate) Predicate {
	return func(object Object) bool {
		for _, p := range predicates {
			if !p(object) {
				return false
			}
		}
		return true
	}
}

func Or(predicates ...Predicate) Predicate {
	return func(object Object) bool {
		for _, p := range predicates {
			if p(object) {
				return true
			}
		}
		return false
	}
}

func Contains[T Object](objects []T, predicate Predicate) bool {
	for _, object := range objects {
		if predicate(object) {
			return true
		}
	}
	return false
}

func Count[T Object](objects []T, predicate Predicate) int {
	count := 0
	for _, object := range objects {
		if predicate(object) {
			count++
		}
	}
	return count
}

func Find[T Object](objects []T, predicate Predicate) (T, bool) {
	var result []T
	for _, object := range objects {
		if predicate(object) {
			result = append(result, object)
		}
	}
	if len(result) == 0 {
		var zero T
		return zero, false
	}
	return result[0], true
}

func FindAll[T Object](objects []T, predicate Predicate) []T {
	var result []T
	for _, object := range objects {
		if predicate(object) {
			result = append(result, object)
		}
	}
	return result
}

func Get[T Object](objects []T, id string) (T, bool) {
	return Find(objects, func(obj Object) bool {
		return obj.ID() == id
	})
}

func GetAll[T Object](objects []T) []T {
	return slices.Clone(objects)
}
