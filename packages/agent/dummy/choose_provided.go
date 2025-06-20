package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

// ChooseProvided is a dummy agent that lets unit tests specify the choice
// that the player will make. This agent can be used in unit tests where the
// test is expecting a player to make a specific choice, without needing to
// implement the full logic of a player agent.
// A choice will be consumed from the provided IDs matching the prompt type,
// each time the agent is asked to choose a choice from a prompt.
// It will always choose the first match when the choice is a ChooseOne prompt
// or all matches up to the max for a ChooseMany prompt.
// When no more choices are available, it will return an empty choice
// for ChooseOne and ChooseMany prompts, and a zero value for ChooseNumber prompts.
// If SkipOptional is true, it will skip optional choices and return an empty choice
// for ChooseOne and ChooseMany prompts, and a zero value for ChooseNumber prompts.
type ChooseProvided struct {
	playerID      string
	oneChoiceIDs  []string
	manyChoiceIDs [][]string
	numbersChosen []int
	skipOptional  bool
}

type ChooseProvidedConfig struct {
	OneChoiceIDs  []string
	ManyChoiceIDs [][]string
	NumbersChosen []int
	SkipOptional  bool
}

func NewChooseProvided(playerID string, config ChooseProvidedConfig) *ChooseProvided {
	agent := ChooseProvided{
		playerID:      playerID,
		oneChoiceIDs:  config.OneChoiceIDs,
		manyChoiceIDs: config.ManyChoiceIDs,
		numbersChosen: config.NumbersChosen,
		skipOptional:  config.SkipOptional,
	}
	return &agent
}

func (a *ChooseProvided) PlayerID() string {
	return a.playerID
}

func (a *ChooseProvided) GetNextAction(game state.Game) (engine.Action, error) {
	return action.NewPassPriorityAction(a.playerID), nil
}

func (a *ChooseProvided) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		if prompt.Optional && a.skipOptional {
			return nil, nil
		}
		if len(a.oneChoiceIDs) == 0 {
			return nil, nil
		}
		provideChoiceID, remaining := a.oneChoiceIDs[0], a.oneChoiceIDs[1:]
		a.oneChoiceIDs = remaining
		for _, choice := range opts.Choices {
			if choice.ID() == provideChoiceID {
				return choose.ChooseOneResults{
					Choice: choice,
				}, nil
			}
		}
		return nil, nil
	case choose.ChooseManyOpts:
		if prompt.Optional && a.skipOptional {
			return nil, nil
		}
		if len(a.manyChoiceIDs) == 0 {
			return nil, nil
		}
		providedChoiceIDs, remaining := a.manyChoiceIDs[0], a.manyChoiceIDs[1:]
		a.manyChoiceIDs = remaining
		var chosen []choose.Choice
		for _, choice := range opts.Choices {
			if slices.Contains(providedChoiceIDs, choice.ID()) {
				chosen = append(chosen, choice)
			}
		}
		if len(chosen) > opts.Max {
			chosen = chosen[:opts.Max]
		}
		return choose.ChooseManyResults{
			Choices: chosen,
		}, nil
	case choose.ChooseNumberOpts:
		if prompt.Optional && a.skipOptional {
			return nil, nil
		}
		if len(a.numbersChosen) == 0 {
			return nil, nil
		}
		providedNumber, remaining := a.numbersChosen[0], a.numbersChosen[1:]
		a.numbersChosen = remaining
		return choose.ChooseNumberResults{
			Number: providedNumber,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected choice type %s", opts)
	}
}

func (a *ChooseProvided) ReportState(game state.Game) error {
	return nil
}
