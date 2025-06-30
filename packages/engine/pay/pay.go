package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"fmt"
)

func Cost(someCost cost.Cost, object gob.Object, playerID string) ([]event.GameEvent, error) {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return payCompositeCost(c, object, playerID)
	case cost.DiscardACardCost:
		return payDiscardACardCost(c, playerID)
	case cost.DiscardThisCost:
		return payDiscardThisCost(object, playerID), nil
	case cost.LifeCost:
		return payLifeCost(c, playerID), nil
	case cost.ManaCost:
		return payManaCost(c, playerID), nil
	case cost.TapThisCost:
		return payTapCost(object, playerID), nil
	default:
		return nil, fmt.Errorf("unsupported cost type: %T", c)
	}
}

func payCompositeCost(c cost.CompositeCost, object gob.Object, playerID string) ([]event.GameEvent, error) {
	var events []event.GameEvent
	for _, subCost := range c.Costs() {
		subEvents, err := Cost(subCost, object, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to pay composite cost: %w", err)
		}
		events = append(events, subEvents...)
	}
	return events, nil
}

func payDiscardThisCost(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   object.ID(),
		},
	}
}

func payDiscardACardCost(c cost.DiscardACardCost, playerID string) ([]event.GameEvent, error) {
	if c.TargetID() == "" {
		return nil, fmt.Errorf("DiscardACardCost requires a CardID, got empty string")
	}
	return []event.GameEvent{
		event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   c.TargetID(),
		},
	}, nil
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
