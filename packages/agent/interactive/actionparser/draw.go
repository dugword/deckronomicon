package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type DrawCheatCommand struct {
	Player state.Player
}

func (p *DrawCheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *DrawCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewDrawCheatAction(p.Player), nil
}

func parseDrawCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*DrawCheatCommand, error) {
	return nil, nil
}
