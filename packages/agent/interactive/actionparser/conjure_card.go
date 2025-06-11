package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"errors"
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
	cardName string,
	player state.Player,
) (*ConjureCardCheatCommand, error) {
	if cardName == "" {
		return nil, errors.New("conjure card command requires a card name")
	}
	return &ConjureCardCheatCommand{
		Player:   player,
		CardName: cardName,
	}, nil
}
