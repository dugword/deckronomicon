package add

import "deckronomicon/packages/game/gob"

func Item[T gob.Object](objects []T, objs ...T) (newObjects []T) {
	return append(objects[:], objs...)
}
