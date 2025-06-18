package actionparser

import (
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mtg"
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
	if !mtg.IsMana(manaString) {
		return action.AddManaCheatAction{}, fmt.Errorf("string %q is not a valid mana string", manaString)
	}
	return action.NewAddManaCheatAction(player, manaString), nil
}
