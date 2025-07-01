package auto

import (
	"deckronomicon/packages/agent/auto/strategy"
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"deckronomicon/packages/ui"
	"deckronomicon/packages/view"
	"fmt"
)

type RuleBasedAgent struct {
	// Rules       []Rule
	// ChoiceRules []ChoiceRule
	mode        string // The current mode of the player
	playerID    string
	verbose     bool
	LastAction  string
	strategy    *strategy.Strategy
	uiBuffer    *ui.Buffer // Buffer for UI updates
	interactive bool       // Whether the agent is interactive or not
	stops       []mtg.Step
}

func (a *RuleBasedAgent) ReportState(game state.Game) {
	player := game.GetPlayer(a.playerID)
	opponent := game.GetOpponent(a.playerID)
	a.uiBuffer.Update(
		view.NewGameViewFromState(game),
		view.NewPlayerViewFromState(game, player, a.mode),
		view.NewPlayerViewFromState(game, opponent, ""),
	)
}

func NewRuleBasedAgent(
	strategyFile string,
	playerID string,
	displayFile string,
	interactive bool,
	stops []mtg.Step,
	verbose bool,
) (*RuleBasedAgent, error) {
	agent := RuleBasedAgent{
		mode:        "Setup",
		playerID:    playerID,
		verbose:     verbose,
		LastAction:  "",
		strategy:    nil,
		uiBuffer:    ui.NewBuffer(displayFile),
		interactive: interactive,
		stops:       stops,
	}
	strategy, err := strategy.LoadStrategy(strategyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy: %w", err)
	}
	agent.strategy = strategy
	return &agent, nil
}

func (a *RuleBasedAgent) PlayerID() string {
	return a.playerID
}

func (a *RuleBasedAgent) Choose(prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	if a.interactive {
		var choices []choose.Choice
		switch opts := prompt.ChoiceOpts.(type) {
		case choose.ChooseOneOpts:
			choices = opts.Choices
		case choose.ChooseManyOpts:
			choices = opts.Choices
		case choose.MapChoicesToBucketsOpts:
			choices = opts.Choices
		default:
			return nil, fmt.Errorf("unsupported interactive choice options type: %T", prompt.ChoiceOpts)
		}
		fmt.Println("Interactive choice prompt:", prompt.Message)
		a.uiBuffer.UpdateChoices(prompt.Message, choices)
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI buffer: %w", err)
		}
		a.enterToContinue()
	}
	switch opts := prompt.ChoiceOpts.(type) {
	case choose.ChooseOneOpts:
		if len(opts.Choices) == 0 {
			return choose.ChooseOneResults{}, fmt.Errorf("no choices available")
		}
		return choose.ChooseOneOpts{
			Choices: []choose.Choice{opts.Choices[0]},
		}, nil
	case choose.ChooseManyOpts:
		if len(opts.Choices) == 0 || len(opts.Choices) < opts.Min {
			return choose.ChooseManyResults{}, fmt.Errorf("no choices available")
		}
		var preferreedChoices []choose.Choice
		// TODO: Temporary hack, need to let the strategy specify preferred choices for specific prompts
		if prompt.Source.Name() == "Discard to Hand Size" {
			for _, choice := range opts.Choices {
				if choice.Name() == "Island" || choice.Name() == "Mountain" {
					preferreedChoices = append(preferreedChoices, choice)
				}
			}
			if len(preferreedChoices) >= opts.Min {
				return choose.ChooseManyResults{
					Choices: preferreedChoices[:opts.Min], // Return the minimum number of preferred choices
				}, nil
			}
		}
		return choose.ChooseManyResults{
			Choices: opts.Choices[:opts.Min], // Return the minimum number of choices
		}, nil
	case choose.MapChoicesToBucketsOpts:
		if len(opts.Buckets) == 0 {
			return choose.MapChoicesToBucketsResults{}, fmt.Errorf("no buckets available")
		}
		return choose.MapChoicesToBucketsResults{
			Assignments: map[choose.Bucket][]choose.Choice{
				opts.Buckets[0]: opts.Choices, // Assign all choices to the first bucket
			},
		}, nil
	case choose.ChooseNumberOpts:
		return choose.ChooseNumberResults{
			Number: opts.Min, // For now, just return a fixed number
		}, nil
	default:
		return nil, fmt.Errorf("unsupported choice options type: %T", opts)
	}
}

func (a *RuleBasedAgent) GetNextAction(game state.Game) (engine.Action, error) {
	ctx := evalstate.EvalState{
		Game:     game,
		PlayerID: a.playerID,
		Mode:     a.mode,
	}
	for _, mode := range a.strategy.Modes {
		if mode.Name == a.mode {
			continue // Skip the current mode
		}
		if mode.When.Evaluate(&ctx) {
			a.mode = mode.Name
			break
		}
	}
	var act engine.Action
	for _, rule := range a.strategy.Rules[a.mode] {
		if !rule.When.Evaluate(&ctx) {
			continue
		}
		if a.interactive {
			a.uiBuffer.UpdateMessage([]string{fmt.Sprintf("Matched rule: %s", rule.Name)})
		}
		var err error
		act, err = rule.Then.Resolve(&ctx)
		if err != nil {
			fmt.Printf("Error resolving action for rule %s: %v\n", rule.Name, err)
			a.enterToContinue()
			continue
		}
		break // Stop after the first matching rule
	}
	if a.interactive {
		if err := a.uiBuffer.Render(); err != nil {
			return nil, fmt.Errorf("failed to render UI buffer: %w", err)
		}
	}
	if act == nil {
		if game.ActivePlayerID() == a.playerID {
			a.enterToContinueOnSteps(game.Step())
		}
		return action.NewPassPriorityAction(), nil
	}
	a.enterToContinue()
	return act, nil
}
