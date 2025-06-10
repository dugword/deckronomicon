package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type CheatCommand struct {
	Player state.Player
}

func (p *CheatCommand) IsComplete() bool {
	return p.Player.ID() != ""
}

func (p *CheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewCheatAction(p.Player), nil
}

func parseCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*CheatCommand, error) {
	return nil, nil
}
