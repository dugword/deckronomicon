package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/state"
	"fmt"
)

type AddManaCheatAction struct {
	Player     state.Player
	ManaString string
}

func NewAddManaCheatAction(player state.Player, manaString string) AddManaCheatAction {
	return AddManaCheatAction{
		Player:     player,
		ManaString: manaString,
	}
}

func (a AddManaCheatAction) Name() string {
	return "CHEAT: Add Mana"
}

func (a AddManaCheatAction) PlayerID() string {
	return a.Player.ID()
}

func (a AddManaCheatAction) Description() string {
	return "CHEAT: The active player adds mana."
}

func (a AddManaCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a AddManaCheatAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
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
			Player: a.Player.ID(),
		},
	}
	for color, n := range amount.Colors() {
		events = append(events, event.AddManaEvent{
			PlayerID: a.Player.ID(),
			ManaType: color,
			Amount:   n,
		})
	}
	if amount.Generic() > 0 {
		events = append(events, event.AddManaEvent{
			PlayerID: a.Player.ID(),
			Amount:   amount.Generic(),
			ManaType: mana.Colorless,
		})
	}
	return events, nil
}
