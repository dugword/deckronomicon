package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

func PayCost(someCost cost.Cost, object query.Object, player state.Player) ([]event.GameEvent, error) {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return payCompositeCost(c, object, player)
	case cost.ManaCost:
		return payManaCost(c, object, player)
	case cost.TapCost:
		return payTapCost(object, player)
	default:
		return nil, fmt.Errorf("unsupported cost type: %T", c)
	}
}

func payCompositeCost(c cost.CompositeCost, object query.Object, player state.Player) ([]event.GameEvent, error) {
	// Check if the player can pay all parts of the composite cost
	var events []event.GameEvent
	for _, subCost := range c.Costs() {
		subEvents, err := PayCost(subCost, object, player)
		if err != nil {
			return nil, err
		}
		events = append(events, subEvents...)
	}
	return events, nil
}

func payManaCost(c cost.ManaCost, object query.Object, player state.Player) ([]event.GameEvent, error) {
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

func payTapCost(object query.Object, player state.Player) ([]event.GameEvent, error) {
	// Create an event to tap the permanent
	return []event.GameEvent{
		event.TapPermanentEvent{
			PlayerID:    player.ID(),
			PermanentID: object.ID(),
		},
	}, nil
}
