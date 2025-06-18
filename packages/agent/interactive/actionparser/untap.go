package actionparser

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

func parseUntapCheatCommand(
	idOrName string,
	game state.Game,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (action.UntapCheatAction, error) {
	permanents := game.Battlefield().FindAll(
		query.And(has.Controller(player.ID()), is.Tapped()),
	)
	var permanent gob.Permanent
	var err error
	if idOrName == "" {
		permanent, err = buildUntapCommandByChoice(permanents, player, chooseFunc)
		if err != nil {
			return action.UntapCheatAction{}, fmt.Errorf("failed to choose a permanent to untap: %w", err)
		}
	} else {
		found, ok := query.Find(permanents, query.Or(has.ID(idOrName), has.Name(idOrName)))
		if !ok {
			return action.UntapCheatAction{}, fmt.Errorf("failed to find permanent with ID or name %q: %w", idOrName, err)
		}
		permanent = found
	}
	return action.NewUntapCheatAction(player, permanent), nil
}

func buildUntapCommandByChoice(
	permanents []gob.Permanent,
	player state.Player,
	chooseFunc func(prompt choose.ChoicePrompt) (choose.ChoiceResults, error),
) (gob.Permanent, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a permanent to untap",
		Source:   player,
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(permanents),
		},
	}
	choiceResults, err := chooseFunc(prompt)
	if err != nil {
		return gob.Permanent{}, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return gob.Permanent{}, fmt.Errorf("expected a single choice result")
	}
	permanent, ok := selected.Choice.(gob.Permanent)
	if !ok {
		return gob.Permanent{}, fmt.Errorf("selected choice is not a permanent")
	}
	return permanent, nil
}
