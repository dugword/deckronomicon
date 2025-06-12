package remove

import (
	"deckronomicon/packages/query"
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
	for _, object := range objects {
		if !predicate(object) {
			remaining = append(remaining, object)
			continue
		}
	}
	return remaining, len(remaining) < len(objects)
}
