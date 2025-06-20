package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"fmt"
)

// ChooseNoneAgent is a dummy agent that always declines to make a choice,
// even if the choice is required. This agent can be used in unit tests
// where the test is expecting a player to make a choice but does not
// care about the outcome of that choice or where the failure case is
// being tested.
// This agent will often not return a valid choice, which may cause the game to stall.
// When running simulations, it is recommended to use ChooseMinimumAgent instead.
type ChooseNoneAgent struct {
	playerID string
}

func NewChooseNoneAgent(playerID string) *ChooseNoneAgent {
	agent := ChooseNoneAgent{
		playerID: playerID,
	}
	return &agent
}

func (a *ChooseNoneAgent) PlayerID() string {
	return a.playerID
}

func (a *ChooseNoneAgent) GetNextAction(game state.Game) (engine.Action, error) {
	return action.NewPassPriorityAction(a.playerID), nil
}

func (a *ChooseNoneAgent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		_, ok := prompt.ChoiceOpts.(choose.ChooseOneOpts)
		if !ok {
			return nil, fmt.Errorf("expected ChoiceOneOpts, got %T", prompt.ChoiceOpts)
		}
		return choose.ChooseOneResults{
			Choice: nil,
		}, nil
	case choose.ChooseManyOpts:
		_, ok := prompt.ChoiceOpts.(choose.ChooseManyOpts)
		if !ok {
			return nil, fmt.Errorf("expected ChooseManyOpts, got %T", prompt.ChoiceOpts)
		}
		return choose.ChooseManyResults{
			Choices: nil,
		}, nil
	case choose.ChooseNumberOpts:
		_, ok := prompt.ChoiceOpts.(choose.ChooseNumberOpts)
		if !ok {
			return nil, fmt.Errorf("expected ChooseNumberOpts, got %T", prompt.ChoiceOpts)
		}
		return choose.ChooseNumberResults{
			Number: 0,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected choice type %s", opts)
	}
}

func (a *ChooseNoneAgent) ReportState(game state.Game) error {
	return nil
}
