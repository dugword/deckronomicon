package actionparser

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"errors"
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
	cardName string,
	player state.Player,
) (*FindCardCheatCommand, error) {
	if cardName == "" {
		return nil, errors.New("find card command requires a card name")
	}
	return &FindCardCheatCommand{
		Player:   player,
		CardName: cardName,
	}, nil
}
