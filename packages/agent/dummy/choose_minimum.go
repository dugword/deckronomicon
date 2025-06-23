package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"fmt"
)

// ChooseMinimumAgent is a dummy agent that always chooses the minimum number of options available.
// When a choice is optional, it will not make a choice.
// When a choice is required, it will select the minimum number of options specified by the choice prompt.
// Options will be selected in the order they are presented.
// It can be used as the default player agent for opponents in games where the player is not interactive.
// This agent should always return valid choices and prevent the game from stalling.
type ChooseMinimumAgent struct {
	playerID string
}

func NewChooseMinimumAgent(playerID string) *ChooseMinimumAgent {
	agent := ChooseMinimumAgent{
		playerID: playerID,
	}
	return &agent
}

func (a *ChooseMinimumAgent) PlayerID() string {
	return a.playerID
}

func (a *ChooseMinimumAgent) GetNextAction(game state.Game) (engine.Action, error) {
	return action.NewPassPriorityAction(), nil
}

func (a *ChooseMinimumAgent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		if prompt.Optional {
			return choose.ChooseOneResults{}, nil
		}
		if len(opts.Choices) == 0 {
			return nil, fmt.Errorf("no choices available")
		}
		return choose.ChooseOneResults{Choice: opts.Choices[0]}, nil
	case choose.ChooseManyOpts:
		if prompt.Optional || opts.Min == 0 {
			return choose.ChooseManyResults{}, nil
		}
		if len(opts.Choices) < opts.Min {
			return nil, fmt.Errorf("not enough choices available")
		}
		return choose.ChooseManyResults{Choices: opts.Choices[:opts.Min]}, nil
	case choose.ChooseNumberOpts:
		if prompt.Optional || opts.Min == 0 {
			return choose.ChooseNumberResults{}, nil
		}
		return choose.ChooseNumberResults{Number: opts.Min}, nil
	default:
		return nil, fmt.Errorf("unknown choice options type: %T", opts)
	}
}

func (a *ChooseMinimumAgent) ReportState(game state.Game) {}
