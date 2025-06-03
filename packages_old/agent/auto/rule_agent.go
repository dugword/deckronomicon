package auto

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/strategy"
	"deckronomicon/packages/strategy/evalstate"
	"deckronomicon/packages/ui"
	"fmt"
)

type RuleBasedAgent struct {
	// Rules       []Rule
	// ChoiceRules []ChoiceRule
	player      *player.Player
	verbose     bool
	LastAction  string
	strategy    *strategy.Strategy
	uiBuffer    *ui.Buffer // Buffer for UI updates
	interactive bool       // Whether the agent is interactive or not
}

func (a *RuleBasedAgent) RegisterPlayer(player *player.Player) {
	a.player = player
}

// TODO: This is just for now, but I like being able to walk through the
// rules agent.
func (a *RuleBasedAgent) ReportState(state player.GameState) {
}

func NewRuleBasedAgent(strategyFile string, interactive bool) (*RuleBasedAgent, error) {
	agent := RuleBasedAgent{
		// TODO : make this configurable
		verbose:     true,
		interactive: interactive,
		uiBuffer:    ui.NewBuffer(),
	}
	strategy, err := strategy.LoadStrategy(strategyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy: %w", err)
	}
	agent.strategy = strategy
	return &agent, nil
}

func (a *RuleBasedAgent) ChooseMany(prompt string, source choose.Source, choices []choose.Choice) ([]choose.Choice, error) {
	panic("not yet implemented")
}

func (a *RuleBasedAgent) ChooseOne(prompt string, source choose.Source, choices []choose.Choice) (choose.Choice, error) {
	return choices[0], nil // For now, just return the first choice
}

func (a *RuleBasedAgent) Confirm(prompt string, source choose.Source) (bool, error) {
	return true, nil
}

func (a *RuleBasedAgent) GetNextAction(state player.GameState) (game.Action, error) {
	s, ok := state.(*engine.GameState)
	if !ok {
		return game.Action{}, fmt.Errorf("expected GameState, got %T", state)
	}
	ctx := &evalstate.EvalState{
		State:       s,
		Player:      a.player,
		Definitions: a.strategy.Definitions,
	}
	for _, mode := range a.strategy.Modes {
		if mode.Name == a.player.Mode {
			continue // Skip the current mode
		}
		result, err := mode.When.Evaluate(ctx)
		if err != nil {
			return game.Action{}, fmt.Errorf("failed to evaluate condition for mode %s: %w", mode.Name, err)
		}
		if result {
			a.player.Mode = mode.Name
			break
		}
	}
	var act game.Action
	var matchedRule string
	for _, rule := range a.strategy.Rules[a.player.Mode] {
		result, err := rule.When.Evaluate(ctx)
		if err != nil {
			return game.Action{}, fmt.Errorf("failed to evaluate condition for rule %s: %w", rule.Name, err)
		}
		if !result {
			continue
		}
		matchedRule = rule.Name
		act, err = rule.Then.Resolve(ctx)
		if err != nil {
			return game.Action{}, fmt.Errorf("failed to resolve action for rule %s: %w", rule.Name, err)
		}
		break // Stop after the first matching rule
	}
	a.uiBuffer.UpdateFromState(s, a.player)
	if err := a.uiBuffer.Render(); err != nil {
		return game.Action{}, fmt.Errorf("failed to render UI buffer: %w", err)
	}
	if string(act.Type) == "" {
		fmt.Println("No action matched for player:", a.player.ID())
		enterToContinue()
		return game.Action{Type: game.ActionPass}, nil // No action matched, just pass
	}
	fmt.Println("Matched rule: ", matchedRule)
	fmt.Println("Action chosen for player: ", act.Type)
	enterToContinue()
	return act, nil
}

func (a *RuleBasedAgent) EnterNumber(string, choose.Source) (int, error) {
	// For now, just return a fixed number
	return 0, nil
}
