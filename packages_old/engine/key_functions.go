package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/query/has"
	"fmt"
)

// TODO: This file name sucks, should be something like `game_actions.go`, or
// rules_actions.go, or something like that.
// Are these Actions? I think they are... that whole concept needs some
// overhaul.

func DrawStartingHand(plyer *player.Player, startingHand []string) ([]string, error) {
	var cardsDrawn []string
	for _, cardName := range startingHand {
		if err := plyer.Tutor(has.Name(cardName)); err != nil {
			return nil, fmt.Errorf("failed to tutor card %s: %w", cardName, err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	remainingDraws := plyer.MaxHandSize - plyer.Hand().Size()
	for range remainingDraws {
		cardName, err := plyer.DrawCard()
		if err != nil {
			return nil, fmt.Errorf("failed to draw card: %w", err)
		}
		cardsDrawn = append(cardsDrawn, cardName)
	}
	return cardsDrawn, nil
}

// Mulligan allows the player to mulligan their hand. The player draws 7 new
// cards but then needs to put 1 card back on the bottom of their library for
// each time they've mulliganed.
func Mulligan(state *GameState, player *player.Player, startingHand []string) error {
	var accept bool
	for (player.Mulligans < player.MaxHandSize) || !accept {
		player.Agent.ReportState(state)
		accept, err := player.Agent.Confirm("Keep Hand? (y/n)", nil)
		if err != nil {
			return fmt.Errorf("failed to confirm mulligan: %w", err)
		}
		if accept {
			break
		}
		for _, card := range player.Hand().GetAll() {
			player.BottomCard(card.ID())
		}
		player.ShuffleLibrary()
		DrawStartingHand(player, startingHand)
		player.Mulligans++
	}
	for range player.Mulligans {
		hand := player.Hand().GetAll()
		choices := choose.CreateChoices(hand, player.Hand())
		choice, err := player.Agent.ChooseOne(
			"Choose a card to put back on the bottom of your library",
			// TODO: Handle this as a constant, not sure what "Mulligan" is
			// though - but need to classify it and use that constant.
			choose.NewChoiceSource("Mulligan"),
			choices,
		)
		if err != nil {
			return fmt.Errorf("failed to choose card to put back: %w", err)
		}
		if err := player.BottomCard(choice.ID); err != nil {
			return fmt.Errorf("failed to put card back on bottom: %w", err)
		}
	}
	return nil
}
