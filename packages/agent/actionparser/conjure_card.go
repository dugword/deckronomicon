package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ConjureCardCheatCommand struct {
	PlayerID string
	CardName string
}

func (p *ConjureCardCheatCommand) IsComplete() bool {
	return p.PlayerID != "" && p.CardName != ""
}

func (p *ConjureCardCheatCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewConjureCardCheatAction(p.PlayerID, p.CardName), nil
}

func parseConjureCardCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*ConjureCardCheatCommand, error) {
	return nil, nil
}
