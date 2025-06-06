package take

import (
	"deckronomicon/packages/query"
	"slices"
)

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

func By[T query.Object](objects []T, predicate query.Predicate) (taken T, remaining []T, ok bool) {
	for i, object := range objects {
		if !predicate(object) {
			continue
		}
		remaining = slices.Delete(objects, i, i+1)
		return object, remaining, true
	}
	return taken, objects, false
}

func Top[T query.Object](objects []T) (taken T, remaining []T, ok bool) {
	if len(objects) == 0 {
		return taken, objects, false
	}
	return objects[0], objects[1:], true
}
