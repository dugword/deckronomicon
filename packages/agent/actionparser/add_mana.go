package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type AddManaCheatCommand struct {
	PlayerID string
	Mana     string
}

func (p *AddManaCheatCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Mana != ""
}

func (p *AddManaCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewAddManaCheatAction(p.PlayerID, p.Mana), nil
}

func parseAddManaCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*AddManaCheatCommand, error) {
	return nil, nil
}
