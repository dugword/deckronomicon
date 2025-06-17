package rng

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

func (r *RNG) ShuffleCards(cards []gob.Card) []gob.Card {
	shuffled := slices.Clone(cards)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

func (r *RNG) ShuffleCardsIDs(cards []gob.Card) []string {
	var shuffled []string
	for _, card := range cards {
		shuffled = append(shuffled, card.ID())
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
