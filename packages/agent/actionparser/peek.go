package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type PeekCheatCommand struct {
	PlayerID string
}

func (p *PeekCheatCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *PeekCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewPeekCheatAction(p.PlayerID), nil
}

func parsePeekCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*PeekCheatCommand, error) {
	return &PeekCheatCommand{PlayerID: playerID}, nil
}
