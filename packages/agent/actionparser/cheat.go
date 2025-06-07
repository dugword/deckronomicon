package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type CheatCommand struct {
	PlayerID string
}

func (p *CheatCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *CheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewCheatAction(p.PlayerID), nil
}

func parseCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*CheatCommand, error) {
	return nil, nil
}
