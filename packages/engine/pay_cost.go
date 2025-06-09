package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/state"
	"fmt"
)

func PayCost(someCost cost.Cost, game state.Game, player state.Player) ([]event.GameEvent, error) {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return payCompositeCost(c, game, player)
	case cost.ManaCost:
		return payManaCost(c, game, player)
	case cost.TapCost:
		return payTapCost(c, game, player)
	default:
		return nil, fmt.Errorf("unsupported cost type: %T", c)
	}
}

func payCompositeCost(c cost.CompositeCost, game state.Game, player state.Player) ([]event.GameEvent, error) {
	// Check if the player can pay all parts of the composite cost
	var events []event.GameEvent
	for _, subCost := range c.Costs() {
		subEvents, err := PayCost(subCost, game, player)
		if err != nil {
			return nil, err
		}
		events = append(events, subEvents...)
	}
	return events, nil
}

func payManaCost(c cost.ManaCost, game state.Game, player state.Player) ([]event.GameEvent, error) {
	/*
		// Check if the player has enough mana to pay the cost
		availableMana := player.ManaAvailable()
		for _, mana := range c.Mana() {
			if availableMana[mana.Color] < mana.Amount {
				return nil, fmt.Errorf("not enough mana to pay cost: %s", c)
			}
		}

		// Create an event to deduct the mana
		return []event.GameEvent{
			event.DeductManaEvent{
				PlayerID:   player.ID(),
				ManaAmount: c.Mana(),
			},
		}, nil
	*/
	return nil, nil
}

func payTapCost(c cost.TapCost, game state.Game, player state.Player) ([]event.GameEvent, error) {
	// Create an event to tap the permanent
	fmt.Println("Tapping permanent:", c.Permanent().ID())
	return []event.GameEvent{
		event.NewTapPermanentEvent(c.Permanent().ID()),
	}, nil
}
