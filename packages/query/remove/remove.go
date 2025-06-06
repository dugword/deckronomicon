package remove

import (
	"deckronomicon/packages/query"
	"slices"
)

func AllBy[T query.Object](objects []T, predicate query.Predicate) (remaining []T) {
	for _, object := range objects {
		if !predicate(object) {
			remaining = append(remaining, object)
			continue
		}
	}
	return remaining
}

func By[T query.Object](objects []T, predicate query.Predicate) (remaining []T, ok bool) {
	for i, object := range objects {
		if !predicate(object) {
			continue
		}
		remaining = slices.Delete(objects, i, i+1)
		return remaining, true
	}
	return objects, false
}
