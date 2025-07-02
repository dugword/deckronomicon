package rng

import (
	"math/rand"
	"slices"
)

type RNG struct {
	rng *rand.Rand
}

func NewRNG(seed int64) *RNG {
	return &RNG{
		rng: rand.New(rand.NewSource(seed)),
	}
}

/* // I think I don't need this, use shuffleIDs instead
func (r *RNG) ShuffleCards(cards []*gob.Card) []*gob.Card {
	shuffled := slices.Clone(cards)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
*/

func (r *RNG) ShuffleIDs(ids []string) []string {
	shuffled := slices.Clone(ids)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
