package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Figure out how to avoid stuttering pay.PayCost
func PayCost(someCost cost.Cost, object query.Object, player state.Player) []event.GameEvent {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return payCompositeCost(c, object, player)
	case cost.ManaCost:
		return payManaCost(c, player)
	case cost.TapThisCost:
		return payTapCost(object, player)
	case cost.DiscardThisCost:
		return payDiscardCost(object, player)
	case cost.LifeCost:
		return payLifeCost(c, player)
	default:
		panic(fmt.Errorf("unsupported cost type: %T", c))
	}
}

func payLifeCost(c cost.LifeCost, player state.Player) []event.GameEvent {
	// Create an event to pay the life cost
	return []event.GameEvent{
		event.LoseLifeEvent{
			PlayerID: player.ID(),
			Amount:   c.Amount(),
		},
	}
}

func payCompositeCost(c cost.CompositeCost, object query.Object, player state.Player) []event.GameEvent {
	// Check if the player can pay all parts of the composite cost
	var events []event.GameEvent
	for _, subCost := range c.Costs() {
		subEvents := PayCost(subCost, object, player)
		events = append(events, subEvents...)
	}
	return events
}

func payManaCost(c cost.ManaCost, player state.Player) []event.GameEvent {
	if c.Amount().Total() == 0 {
		return nil
	}
	return []event.GameEvent{event.SpendManaEvent{
		PlayerID:   player.ID(),
		ManaString: c.Amount().ManaString(),
	}}
}

func payTapCost(object query.Object, player state.Player) []event.GameEvent {
	// Create an event to tap the permanent
	return []event.GameEvent{
		event.TapPermanentEvent{
			PlayerID:    player.ID(),
			PermanentID: object.ID(),
		},
	}
}

func payDiscardCost(object query.Object, player state.Player) []event.GameEvent {
	// Create an event to discard the specified object
	return []event.GameEvent{
		event.DiscardCardEvent{
			PlayerID: player.ID(),
			CardID:   object.ID(),
		},
	}
}
