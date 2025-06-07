package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type DiscardCheatCommand struct {
	PlayerID string
	Card     string
}

func (p *DiscardCheatCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Card != ""
}

func (p *DiscardCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewDiscardCheatAction(p.PlayerID, p.Card), nil
}

func parseDiscardCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*DiscardCheatCommand, error) {
	return nil, nil
}
