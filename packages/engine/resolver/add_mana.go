package resolver

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"fmt"
)

func ResolveAddMana(
	game *state.Game,
	playerID string,
	addMana *effect.AddMana,
) (Result, error) {
	amount, err := mana.ParseManaString(addMana.Mana)
	if err != nil {
		return Result{}, fmt.Errorf("failed to parse mana string %q: %w", addMana.Mana, err)
	}
	var events []event.GameEvent
	if amount.Generic() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Generic(),
			Color:    mana.Colorless,
		})
	}
	if amount.Colorless() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Colorless(),
			Color:    mana.Colorless,
		})
	}
	if amount.White() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.White(),
			Color:    mana.White,
		})
	}
	if amount.Blue() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Blue(),
			Color:    mana.Blue,
		})
	}
	if amount.Black() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Black(),
			Color:    mana.Black,
		})
	}
	if amount.Red() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Red(),
			Color:    mana.Red,
		})
	}
	if amount.Green() > 0 {
		events = append(events, &event.AddManaEvent{
			PlayerID: playerID,
			Amount:   amount.Green(),
			Color:    mana.Green,
		})
	}
	return Result{
		Events: events,
	}, nil
}
