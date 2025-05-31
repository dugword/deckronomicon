package engine

import (
	"deckronomicon/packages/game/gob"
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

func (r *RNG) DeckShuffler() func(deck []gob.Card) []gob.Card {
	return func(deck []gob.Card) []gob.Card {
		newDeck := slices.Clone(deck)
		r.rng.Shuffle(len(newDeck), func(i, j int) {
			newDeck[i], newDeck[j] = newDeck[j], newDeck[i]
		})
		return newDeck
	}
}
