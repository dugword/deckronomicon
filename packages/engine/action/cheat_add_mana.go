package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type AddManaCheatAction struct {
	ManaString string
}

func NewAddManaCheatAction(player state.Player, manaString string) AddManaCheatAction {
	return AddManaCheatAction{
		ManaString: manaString,
	}
}

func (a AddManaCheatAction) Name() string {
	return "CHEAT: Add Mana"
}

func (a AddManaCheatAction) Complete(game state.Game, player state.Player, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	if a.ManaString == "" {
		return nil, fmt.Errorf("add mana action is missing mana string")
	}
	if !mtg.IsMana(a.ManaString) {
		return nil, fmt.Errorf("string %q is not a valid mana string", a.ManaString)
	}
	amount, err := mana.ParseManaString(a.ManaString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mana string %q: %w", a.ManaString, err)
	}
	events := []event.GameEvent{
		event.CheatAddManaEvent{
			Player: player.ID(),
		},
	}
	for color, n := range amount.Colors() {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			ManaType: color,
			Amount:   n,
		})
	}
	if amount.Generic() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: player.ID(),
			Amount:   amount.Generic(),
			ManaType: mana.Colorless,
		})
	}
	return events, nil
}
