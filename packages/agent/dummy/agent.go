package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

type Agent struct {
	id      string
	verbose bool
	stops   []mtg.Step
}

func NewAgent(id string, stops []mtg.Step, verbose bool) *Agent {
	agent := Agent{
		id:      id,
		stops:   stops,
		verbose: verbose,
	}
	return &agent
}

func (a *Agent) PlayerID() string {
	return a.id
}

func (a *Agent) GetNextAction(game state.Game) (engine.Action, error) {
	player, ok := game.GetPlayer(a.id)
	if !ok {
		return nil, fmt.Errorf("player %q not found", a.id)
	}
	return action.NewPassPriorityAction(player), nil
}

func (a *Agent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		if opts.Optional {
			return choose.ChooseOneResults{}, nil
		}
		if len(opts.Choices) == 0 {
			return nil, fmt.Errorf("no choices available")
		}
		return choose.ChooseOneResults{Choice: opts.Choices[0]}, nil
	case choose.ChooseManyOpts:
		if opts.Optional || opts.Min == 0 {
			return choose.ChooseManyResults{}, nil
		}
		if len(opts.Choices) < opts.Min {
			return nil, fmt.Errorf("not enough choices available")
		}
		return choose.ChooseManyResults{Choices: opts.Choices[:opts.Min]}, nil
	default:
		return nil, fmt.Errorf("unknown choice options type: %T", opts)
	}
}

func (a *Agent) ReportState(game state.Game) error {
	return nil
}
