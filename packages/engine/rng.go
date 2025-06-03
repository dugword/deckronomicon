package engine

import "math/rand"

type RNG struct {
	r *rand.Rand
}

func NewRNG(seed int64) *RNG {
	return &RNG{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (r *RNG) Int63() int64 {
	return r.r.Int63()
}
