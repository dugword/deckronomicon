package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ViewCommand struct {
	PlayerID string
	Card     string
	Zone     string
}

func (p *ViewCommand) IsComplete() bool {
	return p.PlayerID != "" && p.Card != ""
}

func (p *ViewCommand) Build(game state.Game, playerID string) (engine.Action, error) {
	return engine.NewViewAction(p.PlayerID, p.Zone, p.Card), nil
}

func parseViewCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	playerID string,
) (*ViewCommand, error) {
	return nil, nil
}
