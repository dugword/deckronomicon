package actionparser

import (
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
	"fmt"
)

func parseAddManaCheatCommand(
	manaString string,
	player state.Player,
) (action.AddManaCheatAction, error) {
	if manaString == "" {
		return action.AddManaCheatAction{}, fmt.Errorf("add mana command requires a mana string")
	}
	if _, err := mana.ParseManaString(manaString); err != nil {
		return action.AddManaCheatAction{}, fmt.Errorf("invalid mana string %q: %w", manaString, err)
	}
	return action.NewAddManaCheatAction(player, manaString), nil
}
