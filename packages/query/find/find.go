package find

import (
	"deckronomicon/packages/query"
)

/*
	func By(objects []object.Object, filter predicate.Predicate) []object.Object {
		var result []object.Object
		for _, object := range objects {
			if filter(object) {
				result = append(result, object)
			}
		}
		return result
	}
*/

func By[T query.Object](objects []T, predicate query.Predicate) []T {
	var result []T
	for _, object := range objects {
		if predicate(object) {
			result = append(result, object)
		}
	}
	return result
}

/*
func FindInZoneBy(zone Zone, filter FilterFunc) []game.Object {
	objects := zone.GetAll()
	return FindBy(objects, filter)
}

func FindBy(objects []game.Object, filter FilterFunc) []game.Object {
	var result []game.Object
	for _, object := range objects {
		if filter(object) {
			result = append(result, object)
		}
	}
	return result
}

func FindFirstInZoneBy(zone Zone, filter FilterFunc) (game.Object, error) {
	objects := zone.GetAll()
	return FindFirstBy(objects, filter)
}

func FindFirstBy(objects []game.Object, filter FilterFunc) (game.Object, error) {
	for _, object := range objects {
		if filter(object) {
			return object, nil
		}
	}
	return nil, ErrObjectNotFound
}

func TakeBy(objects []game.Object, filter FilterFunc) (taken, remaining []game.Object, err error) {
	var result []game.Object
	for i, object := range objects {
		if filter(object) {
			taken = append(result, object)
			// Remove the object from the original slice
			remaining = append(objects[:i], objects[i+1:]...)
			return taken, remaining, nil
		}
	}
	return nil, nil, ErrObjectNotFound
}

func TakeFirstBy(objects []game.Object, filter FilterFunc) (taken game.Object, remaining []game.Object, err error) {
	for i, object := range objects {
		if filter(object) {
			// Remove the object from the original slice
			remaining = append(objects[:i], objects[i+1:]...)
			return object, remaining, nil
		}
	}
	return nil, nil, ErrObjectNotFound
}

func And(filters ...FilterFunc) FilterFunc {
	return func(object game.Object) bool {
		for _, f := range filters {
			if !f(object) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...FilterFunc) FilterFunc {
	return func(object game.Object) bool {
		for _, f := range filters {
			if f(object) {
				return true
			}
		}
		return false
	}
}

// Helpers for common card attributes

func HadID(id string) FilterFunc {
	return func(object game.Object) bool {
		return object.ID() == id
	}
}

func HasColor(color game.Color) FilterFunc {
	return func(object game.Object) bool {
		return object.HasColor(color)
	}
}

func HasCardType(cardType game.CardType) FilterFunc {
	return func(object game.Object) bool {
		return object.HasCardType(cardType)
	}
}

func HasID(id string) FilterFunc {
	return func(object game.Object) bool {
		return object.ID() == id
	}
}

func HasManaValue(manaValue int) FilterFunc {
	return func(object game.Object) bool {
		return object.ManaValue() == manaValue
	}
}

func HasName(name string) FilterFunc {
	return func(object game.Object) bool {
		return object.Name() == name
	}
}

func HasStaticAbility(ability string) FilterFunc {
	return func(object game.Object) bool {
		return object.HasStaticAbility(ability)
	}
}

func HasStaticAbilityModifier(keyword string, modifier EffectTag) FilterFunc {
	return func(object game.Object) bool {
		for _, ability := range object.StaticAbilities() {
			if ability.ID == keyword {
				for _, mod := range ability.Modifiers {
					if mod == modifier {
						return true
					}
				}
			}
		}
		return false
	}
}

func HasSubtype(subtype game.Subtype) FilterFunc {
	return func(object game.Object) bool {
		return object.HasSubtype(subtype)
	}
}

func IsTapped() FilterFunc {
	return func(object game.Object) bool {
		if permanent, ok := object.(*Permanent); ok {
			return permanent.IsTapped()
		}
		return false
	}
}
*/
