package game

type FilterFunc func(GameObject) bool

func FindInZoneBy(zone Zone, filter FilterFunc) []GameObject {
	objects := zone.GetAll()
	return FindBy(objects, filter)
}

func FindBy(objects []GameObject, filter FilterFunc) []GameObject {
	var result []GameObject
	for _, object := range objects {
		if filter(object) {
			result = append(result, object)
		}
	}
	return result
}

func FindFirstInZoneBy(zone Zone, filter FilterFunc) (GameObject, error) {
	objects := zone.GetAll()
	return FindFirstBy(objects, filter)
}

func FindFirstBy(objects []GameObject, filter FilterFunc) (GameObject, error) {
	for _, object := range objects {
		if filter(object) {
			return object, nil
		}
	}
	return nil, ErrObjectNotFound
}

func TakeBy(objects []GameObject, filter FilterFunc) (taken, remaining []GameObject, err error) {
	var result []GameObject
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

func TakeFirstBy(objects []GameObject, filter FilterFunc) (taken GameObject, remaining []GameObject, err error) {
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
	return func(object GameObject) bool {
		for _, f := range filters {
			if !f(object) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...FilterFunc) FilterFunc {
	return func(object GameObject) bool {
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
	return func(object GameObject) bool {
		return object.ID() == id
	}
}

func HasColor(color Color) FilterFunc {
	return func(object GameObject) bool {
		return object.HasColor(color)
	}
}

func HasCardType(cardType CardType) FilterFunc {
	return func(object GameObject) bool {
		return object.HasCardType(cardType)
	}
}

func HasID(id string) FilterFunc {
	return func(object GameObject) bool {
		return object.ID() == id
	}
}

func HasManaValue(manaValue int) FilterFunc {
	return func(object GameObject) bool {
		return object.ManaValue() == manaValue
	}
}

func HasName(name string) FilterFunc {
	return func(object GameObject) bool {
		return object.Name() == name
	}
}

func HasStaticAbility(ability string) FilterFunc {
	return func(object GameObject) bool {
		return object.HasStaticAbility(ability)
	}
}

func HasStaticAbilityModifier(keyword string, modifier EffectTag) FilterFunc {
	return func(object GameObject) bool {
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

func HasSubtype(subtype Subtype) FilterFunc {
	return func(object GameObject) bool {
		return object.HasSubtype(subtype)
	}
}

func IsTapped() FilterFunc {
	return func(object GameObject) bool {
		if permanent, ok := object.(*Permanent); ok {
			return permanent.IsTapped()
		}
		return false
	}
}
