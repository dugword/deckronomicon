package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"fmt"
)

func Cost(someCost cost.Cost, object gob.Object, playerID string) []event.GameEvent {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return payCompositeCost(c, object, playerID)
	case cost.ManaCost:
		return payManaCost(c, playerID)
	case cost.TapThisCost:
		return payTapCost(object, playerID)
	case cost.DiscardThisCost:
		return payDiscardCost(object, playerID)
	case cost.LifeCost:
		return payLifeCost(c, playerID)
	default:
		panic(fmt.Errorf("unsupported cost type: %T", c))
	}
}

func payCompositeCost(c cost.CompositeCost, object gob.Object, playerID string) []event.GameEvent {
	var events []event.GameEvent
	for _, subCost := range c.Costs() {
		subEvents := Cost(subCost, object, playerID)
		events = append(events, subEvents...)
	}
	return events
}

func payLifeCost(c cost.LifeCost, playerID string) []event.GameEvent {
	return []event.GameEvent{
		event.LoseLifeEvent{
			PlayerID: playerID,
			Amount:   c.Amount(),
		},
	}
}

func payManaCost(c cost.ManaCost, playerID string) []event.GameEvent {
	if c.Amount().Total() == 0 {
		return nil
	}
	return []event.GameEvent{event.SpendManaEvent{
		PlayerID:   playerID,
		ManaString: c.Amount().ManaString(),
	}}
}

func payTapCost(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		event.TapPermanentEvent{
			PlayerID:    playerID,
			PermanentID: object.ID(),
		},
	}
}

func payDiscardCost(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   object.ID(),
		},
	}
}
