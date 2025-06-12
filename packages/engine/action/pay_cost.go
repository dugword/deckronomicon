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
	case cost.TapThisCost:
		return payTapCost(object, player)
	case cost.DiscardThisCost:
		return payDiscardCost(object, player)
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
	// Check if the player has enough mana to pay the cost
	if _, err := player.ManaPool().WithSpentFromManaAmount(c.Amount()); err != nil {
		return nil, fmt.Errorf("failed to deduct mana from pool: %w", err)
	}
	return []event.GameEvent{event.SpendManaEvent{
		PlayerID:   player.ID(),
		ManaString: c.Amount().ManaString(),
	}}, nil
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

func payDiscardCost(object query.Object, player state.Player) ([]event.GameEvent, error) {
	// Create an event to discard the specified object
	return []event.GameEvent{
		event.DiscardCardEvent{
			PlayerID: player.ID(),
			CardID:   object.ID(),
		},
	}, nil
}
