package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
)

type ConjureCardCheatCommand struct {
	Player   state.Player
	CardName string
}

func (p *ConjureCardCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.CardName != ""
}

func (p *ConjureCardCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return engine.NewConjureCardCheatAction(p.Player, p.CardName), nil
}

func parseConjureCardCheatCommand(
	command string,
	args []string,
	getChoices func(prompt choose.ChoicePrompt) ([]choose.Choice, error),
	game state.Game,
	player state.Player,
) (*ConjureCardCheatCommand, error) {
	return nil, nil
}
