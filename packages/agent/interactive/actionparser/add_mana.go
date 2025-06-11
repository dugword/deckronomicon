package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type AddManaCheatCommand struct {
	Player     state.Player
	ManaString string
}

func (p *AddManaCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.ManaString != ""
}

func (p *AddManaCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewAddManaCheatAction(p.Player, p.ManaString), nil
}

func parseAddManaCheatCommand(
	manaString string,
	player state.Player,
) (*AddManaCheatCommand, error) {
	if manaString == "" {
		return nil, fmt.Errorf("add mana command requires a mana string")
	}
	if !mtg.IsMana(manaString) {
		return nil, fmt.Errorf("string %q is not a valid mana string", manaString)
	}
	return &AddManaCheatCommand{
		Player:     player,
		ManaString: manaString,
	}, nil
}
