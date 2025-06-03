package take

import "deckronomicon/packages/query"

func AllBy[T query.Object](objects []T, predicate query.Predicate) (taken []T, remaining []T) {
	for _, object := range objects {
		if !predicate(object) {
			remaining = append(remaining, object)
			continue
		}
		taken = append(taken, object)
	}
	return taken, remaining
}

func By[T query.Object](objects []T, predicate query.Predicate) (taken T, remaining []T, err error) {
	for i, object := range objects {
		if !predicate(object) {
			continue
		}
		remaining = append(objects[:1], objects[i+1:]...)
		return object, remaining, nil
	}
	return taken, objects, query.ErrNotFound
}
