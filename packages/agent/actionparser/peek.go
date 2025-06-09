package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type PeekCheatCommand struct {
	Player state.Player
}

func (p *PeekCheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *PeekCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return engine.NewPeekCheatAction(player), nil
}

func parsePeekCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*PeekCheatCommand, error) {
	return &PeekCheatCommand{Player: player}, nil
}
