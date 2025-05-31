package add

import (
	"deckronomicon/packages/query"
)

func Item[T query.Object](objects []T, objs ...T) (newObjects []T) {
	return append(objects[:], objs...)
}
