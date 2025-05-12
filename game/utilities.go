package game

import (
	"fmt"
	"math/rand/v2"
)

// TODO: I hate this file, I shouldn't have a utilities file. These should be better organized some how...

// TODO: Maybe have an abstraction like an interface so I can search hands and graveyards too...
func FindPermanentsWithActivatableAbilities(battlefield []*Permanent) []Choice {
	var found []Choice
	for i, permanent := range battlefield {
		if len(permanent.ActivatedAbilities) > 0 {
			found = append(found, Choice{Name: permanent.Name, Index: i})
		}
	}
	return found
}

func UntapPermanents(battlefield []*Permanent) {
	for _, permanent := range battlefield {
		permanent.Tapped = false
	}
}

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

// Does it matter if there are duplicates? Do I need to track ID? Probably someday, but not now
func removeCard(slice []*Card, target *Card) []*Card {
	for i, c := range slice {
		if c == target {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func findUntappedPermanent(permanents []*Permanent, name string) *Permanent {
	for _, p := range permanents {
		if p.Name == name && !p.Tapped {
			return p
		}
	}
	return nil
}

func findCard(cards []*Card, name string) *Card {
	for _, c := range cards {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func getCardOptions(cards []*Card) []Choice {
	var options []Choice
	for i, card := range cards {
		options = append(options, Choice{Name: card.Name, Index: i})
	}
	return options
}

func GetPotentialMana(state *GameState) map[string]int {
	tempGameState := GameState{
		ManaPool: make(ManaPool),
	}

	for _, permanent := range state.Battlefield {
		if permanent.Tapped {
			continue
		}
		for _, ability := range permanent.ActivatedAbilities {
			if ability.IsManaAbility {
				ability.Effect(&tempGameState, nil) // pass mock resolver
			}
		}
	}
	return tempGameState.ManaPool
}
