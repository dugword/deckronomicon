package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ActivateAbilityCommand struct {
	PlayerID string
}

func (p *ActivateAbilityCommand) IsComplete() bool {
	return p.PlayerID != ""
}

func (p *ActivateAbilityCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewActivateAbilityAction(p.PlayerID), nil
}

func parseActivateAbilityCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*ActivateAbilityCommand, error) {
	return nil, nil
}
