package game

import (
	"fmt"
	"math/rand/v2"
)

// TODO: I hate this file, I shouldn't have a utilities file. These should be better organized some how...

// Card utilities,
func shuffleCards(cards []*Card) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func takeNCards(cards []*Card, n int) (taken, remaining []*Card, err error) {
	if len(cards) < n {
		return nil, nil, fmt.Errorf("not enough cards")
	}
	taken, remaining = cards[:n], cards[n:]
	return taken, remaining, nil
}

// Peek returns the top N cards without modifying the pile.
func peek(pile []*Card, n int) []*Card {
	if n > len(pile) {
		n = len(pile)
	}
	return pile[:n]
}

func putOnTop(pile []*Card, cards ...*Card) []*Card {
	return append(cards, pile...)
}

func putOnBottom(pile []*Card, cards ...*Card) []*Card {
	return append(pile, cards...)
}

func GetPotentialMana(state *GameState) map[string]int {
	tempGameState := GameState{
		ManaPool: make(ManaPool),
	}

	for _, permanent := range state.Battlefield.permanents {
		if permanent.IsTapped() {
			continue
		}
		for _, ability := range permanent.ActivatedAbilities() {
			if ability.IsManaAbility {
				ability.Effect(&tempGameState, nil) // pass mock resolver
			}
		}
	}
	return tempGameState.ManaPool
}
