package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"fmt"
)

type AddManaCheatAction struct {
	ManaString string
}

func NewAddManaCheatAction(manaString string) AddManaCheatAction {
	return AddManaCheatAction{
		ManaString: manaString,
	}
}

func (a AddManaCheatAction) Name() string {
	return "CHEAT: Add Mana"
}

func (a AddManaCheatAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	if a.ManaString == "" {
		return nil, fmt.Errorf("add mana action is missing mana string")
	}
	amount, err := mana.ParseManaString(a.ManaString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana string %q: %w", a.ManaString, err)
	}
	events := []event.GameEvent{
		&event.CheatAddManaEvent{
			Player: playerID,
		},
	}
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
	return events, nil
}
