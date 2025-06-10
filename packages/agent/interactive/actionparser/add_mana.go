package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type AddManaCheatCommand struct {
	Player state.Player
	Mana   string
}

func (p *AddManaCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Mana != ""
}

func (p *AddManaCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewAddManaCheatAction(p.Player, p.Mana), nil
}

func parseAddManaCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*AddManaCheatCommand, error) {
	return nil, nil
}
