package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type UntapCheatCommand struct {
	Player    state.Player
	Permanent gob.Permanent
}

func (p *UntapCheatCommand) IsComplete() bool {
	return p.Player.ID() != "" && p.Permanent.ID() != ""
}

func (p *UntapCheatCommand) Build(game state.Game, player state.Player) (engine.Action, error) {
	return action.NewUntapCheatAction(p.Player, p.Permanent), nil
}

func parseUntapCheatCommand(
	idOrName string,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	game state.Game,
	player state.Player,
) (*UntapCheatCommand, error) {
	permanents := game.Battlefield().FindAll(
		query.And(has.Controller(player.ID()), is.Tapped()),
	)
	if idOrName == "" {
		return buildUntapCommandByChoice(permanents, chooseFunc, player)
	}
	return buildUntapCommandByIDOrName(permanents, idOrName, player)
}

func buildUntapCommandByChoice(
	permanents []gob.Permanent,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
	player state.Player,
) (*UntapCheatCommand, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a permanent to untap",
		Source:   CommandSource{"Untap a permanent"},
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(permanents),
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
	permanent, ok := selected.Choice.(gob.Permanent)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a permanent")
	}
	return &UntapCheatCommand{
		Player:    player,
		Permanent: permanent,
	}, nil
}

func buildUntapCommandByIDOrName(
	permanents []gob.Permanent,
	idOrName string,
	player state.Player,
) (*UntapCheatCommand, error) {
	if permanent, ok := query.Find(permanents, query.Or(has.ID(idOrName), has.Name(idOrName))); ok {
		return &UntapCheatCommand{
			Player:    player,
			Permanent: permanent,
		}, nil
	}
	return nil, fmt.Errorf("no permanent found with ID or name %q", idOrName)
}
