package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"fmt"
)

func Cost(someCost cost.Cost, object gob.Object, playerID string) ([]event.GameEvent, error) {
	switch c := someCost.(type) {
	case cost.Composite:
		return payComposite(c, object, playerID)
	case cost.DiscardACard:
		return payDiscardACard(c, playerID)
	case cost.DiscardThis:
		return payDiscardThis(object, playerID), nil
	case cost.Life:
		return payLife(c, playerID), nil
	case cost.Mana:
		return payMana(c, playerID), nil
	case cost.SacrificeThis:
		return paySacrificeThisCost(object, playerID), nil
	case cost.TapThis:
		return payTapCost(object, playerID), nil
	default:
		return nil, fmt.Errorf("unsupported cost type: %T", c)
	}
}

func payComposite(c cost.Composite, object gob.Object, playerID string) ([]event.GameEvent, error) {
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

func payDiscardThis(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		&event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   object.ID(),
		},
	}
}

func payDiscardACard(c cost.DiscardACard, playerID string) ([]event.GameEvent, error) {
	if c.Target().ID == "" {
		return nil, fmt.Errorf("DiscardACard requires a CardID, got empty string")
	}
	return []event.GameEvent{
		&event.DiscardCardEvent{
			PlayerID: playerID,
			CardID:   c.Target().ID,
		},
	}, nil
}

func payLife(c cost.Life, playerID string) []event.GameEvent {
	return []event.GameEvent{
		&event.LoseLifeEvent{
			PlayerID: playerID,
			Amount:   c.Amount(),
		},
	}
}

func payMana(c cost.Mana, playerID string) []event.GameEvent {
	if c.Amount().Total() == 0 {
		return nil
	}
	return []event.GameEvent{&event.SpendManaEvent{
		PlayerID:   playerID,
		ManaString: c.Amount().ManaString(),
	}}
}

func paySacrificeThisCost(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		&event.SacrificePermanentEvent{
			PlayerID:    playerID,
			PermanentID: object.ID(),
		},
	}
}

func payTapCost(object gob.Object, playerID string) []event.GameEvent {
	return []event.GameEvent{
		&event.TapPermanentEvent{
			PlayerID:    playerID,
			PermanentID: object.ID(),
		},
	}
}
