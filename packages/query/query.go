package query

import (
	"deckronomicon/packages/game/mtg"
	"slices"
)

// TODO: Have some way to stringify these interfaces for debugging purposes
// query.Describe = "match (Creature or Enchantment) and (Blue or Red) and (ManaValue 3 or ManaValue 4)"

// Cards, Permanents, Sp
type Object interface {
	Controller() string
	Owner() string
	Description() string
	ID() string
	Match(Predicate) bool
	Name() string
}

// Cards, Permanents, and Spells are all CardObjects
type CardObject interface {
	Object
	CardTypes() []mtg.CardType
	Colors() mtg.Colors
	StaticKeywords() []mtg.StaticKeyword
	Subtypes() []mtg.Subtype
	Supertypes() []mtg.Supertype
	ManaValue() int
}

type StackObject interface {
	Object
	SourceID() string
}

// TODO: I don't like this, but it is a quick way to convert
// a slice of objects to a slice of Object
// Need to update the strategy parser to not require this
func NewQueryObjects[T Object](objects []T) []Object {
	var queryObjects []Object
	for _, obj := range objects {
		queryObjects = append(queryObjects, obj)
	}
	return queryObjects
}

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

func GetN[T Object](objects []T, n int) []T {
	if n <= 0 || n > len(objects) {
		return objects
	}
	return objects[:n]
}

func GetTop[T Object](objects []T) (T, bool) {
	if len(objects) == 0 {
		var zero T
		return zero, false
	}
	return objects[0], true
}
