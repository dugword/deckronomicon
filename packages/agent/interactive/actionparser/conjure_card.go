package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
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
	return action.NewConjureCardCheatAction(p.Player, p.CardName), nil
}

func parseConjureCardCheatCommand(
	arg string,
	game state.Game,
	player state.Player,
) (*ConjureCardCheatCommand, error) {
	return nil, nil
}
