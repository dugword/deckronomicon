package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"

	"errors"
	"fmt"
)

func AddManaEffectHandler(
	game state.Game,
	player state.Player,
	source query.Object,
	modifiers []gob.Tag,
) ([]event.GameEvent, error) {
	var manaString string
	for _, modifier := range modifiers {
		if modifier.Key == "Mana" {
			manaString = modifier.Value
			break
		}
	}
	if manaString == "" {
		return nil, errors.New("effect 'AddMana' requires 'Mana' modifier")
	}
	amount, err := mana.ParseManaString(manaString)
	if err != nil {
		return nil, fmt.Errorf("invalid mana amount '%s': %w", manaString, err)
	}

	// Think through how to best handle this and how the events will be represented in JSON.
	// I could have the mana struct pretty print to a string like "2{R}{G}" or and then reparse it when I apply the event.
	// Or I could generate multiple events for each color of mana like I am doing here.
	var events []event.GameEvent
	for color, amount := range amount.Colors() {
		if amount <= 0 {
			continue // Skip colors with no mana
		}
		events = append(events, event.NewAddManaEvent(
			player.ID(),
			color,
			amount,
		),
		)
	}
	return events, nil
}
