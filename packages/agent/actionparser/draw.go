package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type DrawCheatCommand struct {
	PlayerID string
}

func (p *DrawCheatCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *DrawCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewDrawCheatAction(p.PlayerID), nil
}

func parseDrawCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*DrawCheatCommand, error) {
	return nil, nil
}
