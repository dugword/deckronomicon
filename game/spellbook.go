package game

import "fmt"

// TODO: Rename this file, maybe make it a package

// Abilities

var AbilityIsland = ActivatedAbility{
	Cost: TapCost{},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		g.ManaPool["U"]++
	},
	Description: "{T}: Add {U}.",
	// Update this to just check it via tags, maybe a helper method
	IsManaAbility: true,
	Tags: []AbilityTag{
		{Key: "ManaSource", Value: "U"},
	},
}

var AbilitySwamp = ActivatedAbility{
	Cost: TapCost{},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		g.ManaPool["B"]++
	},
	Description:   "{T}: Add {B}.",
	IsManaAbility: true,
	Tags: []AbilityTag{
		{Key: "ManaSource", Value: "B"},
	},
}

var AbilityMountain = ActivatedAbility{
	Cost: TapCost{},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		g.ManaPool["R"]++
	},
	Description:   "{T}: Add {R}.",
	IsManaAbility: true,
	Tags: []AbilityTag{
		{Key: "ManaSource", Value: "R"},
	},
}

var AbilityForest = ActivatedAbility{
	Cost: TapCost{},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		g.ManaPool["G"]++
	},
	Description:   "{T}: Add {G}.",
	IsManaAbility: true,
	Tags: []AbilityTag{
		{Key: "ManaSource", Value: "G"},
	},
}

var AbilityPlains = ActivatedAbility{
	Cost: TapCost{},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		g.ManaPool["W"]++
	},
	Description:   "{T}: Add {W}.",
	IsManaAbility: true,
	Tags: []AbilityTag{
		{Key: "ManaSource", Value: "W"},
	},
}

var AbilityScry2 = SpellAbility{
	Description: "Scry 2",
	Effect: func(g *GameState, resolver ChoiceResolver) {
		// TODO: Don't access cards directly
		cards := Scry(2, g.Deck.cards, resolver)
		g.Deck.cards = cards
	},
}

// TODO Handle errors
func Scry(n int, cards []*Card, resolver ChoiceResolver) []*Card {
	if n <= 0 || len(cards) == 0 {
		return cards
	}
	if n > len(cards) {
		n = len(cards)
	}
	// TODO: Handle this error
	taken, remaining, _ := takeNCards(cards, n)
	used := make([]bool, len(taken))
	for range len(taken) {
		// Build option list from unplaced cards
		var choices []Choice
		for index, card := range taken {
			if !used[index] {
				choices = append(choices, Choice{
					Name:  card.Name(),
					Index: index,
				})
			}
		}
		chosen := resolver.ChooseOne("Choose a card to place", choices)
		chosenCard := taken[chosen.Index]
		used[chosen.Index] = true
		topBottomchoices := []Choice{
			{Name: "Top"},
			{Name: "Bottom"},
		}
		placement := resolver.ChooseOne(
			fmt.Sprintf("Place %s on top or bottom of your library?", chosenCard.Name),
			topBottomchoices,
		)
		if placement.Index == 0 {
			remaining = putOnTop(remaining, chosenCard)
		} else {
			remaining = putOnBottom(remaining, chosenCard)
		}
	}
	return remaining
}

var AbilityBrainstorm = SpellAbility{
	Description: "Draw three cards, then put two cards from your hand on top of your library in any order.",
	Effect: func(g *GameState, resolver ChoiceResolver) {
		/* // TODO: Figure out a generic draw that's not a draw action maybe
		g.DrawCards(3)
		hand, deck := PutNBackOnTop(2, g.Hand, g.Deck, resolver)
		g.Deck = deck
		g.Hand = hand
		*/
	},
}

// TODO Revist this
func PutNBackOnTop(n int, hand *Hand, deck []*Card, resolver ChoiceResolver) (newHand, newDeck []*Card) {
	for range n {
		choices := hand.CardChoices()
		choice := resolver.ChooseOne("Which card to put back on top", choices)
		card := hand.GetCard(choice.Index)
		if card != nil {
			card := hand.GetCard(choice.Index)
			if card != nil {
				hand.RemoveCard(card)
				deck = putOnTop(deck, card)
			}
		}
	}
	return hand.cards, deck
}

/*
var AbilityCandyTrailSacrifice = ActivatedAbility{
	Cost: CompositeCost{
		Costs: []Cost{
			ManaCost{Generic: 1},
			SacrificeSelfCost{},
		},
	},
	Effect: func(g *GameState, resolver ChoiceResolver) {
		DrawCard(g)
	},
	Description: "{1}, Sacrifice: Draw a card.",
}
*/
