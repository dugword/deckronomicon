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

type DiscardCheatCommand struct {
	Player state.Player
	Card   gob.Card
}

func (p *DiscardCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Card.ID() != ""
}

func (p *DiscardCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewDiscardCheatAction(p.Player, p.Card), nil
}

func parseDiscardCheatCommand(
	idOrName string,
	choose func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (*DiscardCheatCommand, error) {
	cards := player.Hand().GetAll()
	if idOrName == "" {
		return buildDiscardCommandByChoice(cards, choose, player)
	}
	return buildDiscardCommandByIDOrName(cards, idOrName, player)
}

func buildDiscardCommandByChoice(
	cards []gob.Card,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player,
) (*DiscardCheatCommand, error) {
	prompt := choose.ChoicePrompt{
		Message: "Choose a card to discard",
		Source:  CommandSource{"Discard a card"},
		ChoiceOpts: choose.ChooseOneOpts{
			Choices:  choose.NewChoices(cards),
			Optional: true,
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected a single choice result")
	}
	card, ok := selected.Choice.(gob.Card)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a card in hand")
	}
	return &DiscardCheatCommand{
		Card:   card,
		Player: player,
	}, nil
}

func buildDiscardCommandByIDOrName(
	cards []gob.Card,
	idOrName string,
	player state.Player,
) (*DiscardCheatCommand, error) {
	if card, ok := query.Find(cards, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &DiscardCheatCommand{
			Card:   card,
			Player: player,
		}, nil
	}
	return nil, fmt.Errorf("no land found with ID or name %q", idOrName)
}
