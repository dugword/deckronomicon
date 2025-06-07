package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type FindCardCheatCommand struct {
	PlayerID string
	Card     string
}

func (p *FindCardCheatCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Card != ""
}

func (p *FindCardCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewFindCardCheatAction(p.PlayerID, p.Card), nil
}

func parseFindCardCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*FindCardCheatCommand, error) {
	return nil, nil
}
