package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type FindCardCheatCommand struct {
	Player state.Player
	Card   gob.Card
}

func (p *FindCardCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Card.ID() != ""
}

func (p *FindCardCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewFindCardCheatAction(p.Player, p.Card), nil
}

func parseFindCardCheatCommand(
	idOrName string,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	player state.Player,
) (*FindCardCheatCommand, error) {
	cards := player.Library().GetAll()
	if idOrName == "" {
		return buildFindCardCheatCommandByChoice(cards, chooseOne, player)
	}
	return buildFindCardCheatCommandByIDOrName(cards, idOrName, player)
}

func buildFindCardCheatCommandByChoice(
	cards []gob.Card,
	chooseOne func(prompt choose.ChoicePrompt) (choose.Choice, error),
	player state.Player) (*FindCardCheatCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a card to put into your hand",
		Choices:  choose.NewChoices(cards),
		Source:   CommandSource{"Find a card"},
		Optional: true,
	}
	selected, err := chooseOne(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	card, ok := selected.(gob.Card)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in a zone")
	}
	return &FindCardCheatCommand{
		Card:   card,
		Player: player,
	}, nil
}

func buildFindCardCheatCommandByIDOrName(
	cards []gob.Card,
	idOrName string,
	player state.Player,
) (*FindCardCheatCommand, error) {
	if card, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &FindCardCheatCommand{
			Card:   card,
			Player: player,
		}, nil
	}
	return nil, fmt.Errorf("no card found with ID or name %q", idOrName)
}
