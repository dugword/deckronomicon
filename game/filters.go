package game

type FilterFunc func(GameObject) bool

func FindBy(zone Zone, filter FilterFunc) []GameObject {
	objects := zone.GetAll()
	var result []GameObject
	for _, object := range objects {
		if filter(object) {
			result = append(result, object)
		}
	}
	return result
}

func FindFirstBy(zone Zone, filter FilterFunc) (GameObject, error) {
	objects := zone.GetAll()
	for _, object := range objects {
		if filter(object) {
			return object, nil
		}
	}
	return nil, ErrObjectNotFound
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
