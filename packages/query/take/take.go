package take

import (
	"deckronomicon/packages/query"
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
	for _, object := range objects {
		if predicate(object) {
			taken = object
			continue
		}
		remaining = append(remaining, object)
	}
	if len(remaining) == len(objects) {
		return taken, objects, false
	}
	return taken, remaining, true
}

func Top[T query.Object](objects []T) (taken T, remaining []T, ok bool) {
	if len(objects) == 0 {
		return taken, objects, false
	}
	return objects[0], objects[1:], true
}

func N[T query.Object](objects []T, n int) (taken []T, remaining []T) {
	if n <= 0 || len(objects) == 0 {
		return taken, objects
	}
	if n >= len(objects) {
		return objects, nil
	}
	return objects[:n], objects[n:]
}
