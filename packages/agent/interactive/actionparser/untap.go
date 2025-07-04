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

func parseUntapCheatCommand(
	idOrName string,
	game *state.Game,
	playerID string,
	agent engine.PlayerAgent,
) (action.UntapCheatAction, error) {
	permanents := game.Battlefield().FindAll(
		query.And(has.Controller(playerID), is.Tapped()),
	)
	var permanent *gob.Permanent
	var err error
	if idOrName == "" {
		permanent, err = buildUntapCommandByChoice(game, permanents, playerID, agent)
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
	return action.NewUntapCheatAction(permanent.ID()), nil
}

func buildUntapCommandByChoice(
	game *state.Game,
	permanents []*gob.Permanent,
	playerID string,
	agent engine.PlayerAgent,
) (*gob.Permanent, error) {
	prompt := choose.ChoicePrompt{
		Message:  "Choose a permanent to untap",
		Source:   nil, // TODO: Make this better
		Optional: true,
		ChoiceOpts: choose.ChooseOneOpts{
			Choices: choose.NewChoices(permanents),
		},
	}
	choiceResults, err := agent.Choose(game, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get choices: %w", err)
	}
	selected, ok := choiceResults.(choose.ChooseOneResults)
	if !ok {
		return nil, fmt.Errorf("expected a single choice result")
	}
	permanent, ok := selected.Choice.(*gob.Permanent)
	if !ok {
		return nil, fmt.Errorf("selected choice is not a permanent")
	}
	return permanent, nil
}
