package auto

import (
	"deckronomicon/auto/strategy"
	"deckronomicon/game"
	"deckronomicon/ui"
	"fmt"
)

type RuleBasedAgent struct {
	// Rules       []Rule
	// ChoiceRules []ChoiceRule
	playerID    string
	verbose     bool
	LastAction  string
	strategy    *strategy.Strategy
	uiBuffer    *ui.Buffer // Buffer for UI updates
	interactive bool       // Whether the agent is interactive or not
}

func (a *RuleBasedAgent) PlayerID() string {
	return a.playerID
}

// TODO: This is just for now, but I like being able to walk through the
// rules agent.
func (a *RuleBasedAgent) ReportState(state *game.GameState) {
}

func NewRuleBasedAgent(strategyFile string, playerID string, interactive bool) (*RuleBasedAgent, error) {
	agent := RuleBasedAgent{
		// TODO : make this configurable
		verbose:     true,
		playerID:    playerID,
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

func (a *RuleBasedAgent) ChooseMany(prompt string, source game.ChoiceSource, choices []game.Choice) ([]game.Choice, error) {
	panic("not yet implemented")
}

func (a *RuleBasedAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	return choices[0], nil // For now, just return the first choice
}

func (a *RuleBasedAgent) Confirm(prompt string, source game.ChoiceSource) (bool, error) {
	return true, nil
}

func (a *RuleBasedAgent) GetNextAction(state *game.GameState) (*game.GameAction, error) {
	player, err := state.GetPlayer(a.playerID)
	if err != nil {
		return nil, fmt.Errorf("player not found in state: %s", a.playerID)
	}
	ctx := &strategy.EvaluatorContext{
		State:    state,
		Player:   player,
		Strategy: a.strategy,
	}
	for _, mode := range a.strategy.Modes {
		if mode.Name == player.Mode {
			continue // Skip the current mode
		}
		result, err := mode.When.Evaluate(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate condition for mode %s: %w", mode.Name, err)
		}
		if result {
			player.Mode = mode.Name
			break
		}
	}
	var action *game.GameAction
	var matchedRule string
	for _, rule := range a.strategy.Rules[player.Mode] {
		result, err := rule.When.Evaluate(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate condition for rule %s: %w", rule.Name, err)
		}
		if !result {
			continue
		}
		matchedRule = rule.Name
		action, err = rule.Then.Resolve(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve action for rule %s: %w", rule.Name, err)
		}
		break // Stop after the first matching rule
	}
	a.uiBuffer.UpdateFromState(state, a.playerID)
	if err := a.uiBuffer.Render(); err != nil {
		return nil, fmt.Errorf("failed to render UI buffer: %w", err)
	}
	if action == nil {
		fmt.Println("No action matched for player:", a.playerID)
		enterToContinue()
		return &game.GameAction{Type: game.ActionPass}, nil // No action matched, just pass
	}
	fmt.Println("Matched rule: ", matchedRule)
	fmt.Println("Action chosen for player: ", action.Type)
	enterToContinue()
	return action, nil
}

func (a *RuleBasedAgent) EnterNumber(string, game.ChoiceSource) (int, error) {
	// For now, just return a fixed number
	return 0, nil
}
