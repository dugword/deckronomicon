package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

type FindCardCheatCommand struct {
	Player   state.Player
	CardName string
}

func (p *FindCardCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.CardName != ""
}

func (p *FindCardCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewFindCardCheatAction(p.Player, p.CardName), nil
}

func parseFindCardCheatCommand(
	arg string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	game state.Game,
	player state.Player,
) (*FindCardCheatCommand, error) {
	return nil, nil
}
