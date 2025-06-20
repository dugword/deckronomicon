package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"fmt"
)

// ChooseOneAgent is a dummy agent that always chooses the first option available, even
// if the choice is optional. This agent can be used in unit tests where the test is expecting
// a player to make a choice.
// This agent should always return a valid choice, but some choices may require additional
// resources the player does not have, which may cause the game to stall.
// When running simulations, it is recommended to use ChooseMinimumAgent instead.
type ChooseOneAgent struct {
	playerID string
}

func NewChooseOneAgent(playerID string) *ChooseOneAgent {
	agent := ChooseOneAgent{
		playerID: playerID,
	}
	return &agent
}

func (a *ChooseOneAgent) PlayerID() string {
	return a.playerID
}

func (a *ChooseOneAgent) GetNextAction(game state.Game) (engine.Action, error) {
	return action.NewPassPriorityAction(a.playerID), nil
}

func (a *ChooseOneAgent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		if len(opts.Choices) == 0 {
			return nil, nil
		}
		return choose.ChooseOneResults{
			Choice: opts.Choices[0],
		}, nil
	case choose.ChooseManyOpts:
		if len(opts.Choices) == 0 {
			return nil, nil
		}
		return choose.ChooseManyResults{
			Choices: opts.Choices[:1],
		}, nil
	case choose.ChooseNumberOpts:
		return choose.ChooseNumberResults{
			Number: 1,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected choice type %s", opts)
	}
}

func (a *ChooseOneAgent) ReportState(game state.Game) error {
	return nil
}
