package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type DiscardCheatCommand struct {
	Player state.Player
	Card   string
}

func (p *DiscardCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Card != ""
}

func (p *DiscardCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewDiscardCheatAction(p.Player, p.Card), nil
}

func parseDiscardCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*DiscardCheatCommand, error) {
	return nil, nil
}
