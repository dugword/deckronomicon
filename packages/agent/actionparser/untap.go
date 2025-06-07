package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type UntapCheatCommand struct {
	PlayerID string
	Card     string
}

func (p *UntapCheatCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Card != ""
}

func (p *UntapCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewUntapCheatAction(p.PlayerID, p.Card), nil
}

func parseUntapCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*UntapCheatCommand, error) {
	return nil, nil
}
