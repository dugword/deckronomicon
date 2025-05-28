package auto

import (
	"deckronomicon/game"
	"fmt"
)

type EvaluatorContext struct {
	state    *game.GameState
	player   *game.Player
	strategy *Strategy
}

type RuleBasedAgent struct {
	// Rules       []Rule
	// ChoiceRules []ChoiceRule
	playerID   string
	verbose    bool
	LastAction string
	strategy   *Strategy
}

type Rule struct {
	Name        string         `json:"Name"`
	Description string         `json:"Description"`
	RawWhen     map[string]any `json:"When"`
	When        ConditionNode  `json:"-"`
	RawThen     map[string]any `json:"Then"`
	Then        ActionNode     `json:"-"`
}

type Strategy struct {
	Name        string              `json:"Name,omitempty"`
	Description string              `json:"Description,omitempty"`
	Definitions map[string][]string `json:"Definitions,omitempty"`
	Modes       []Rule              `json:"Modes,omitempty"`
	Rules       map[string][]Rule   `json:"Rules,omitempty"`
}

func (a *RuleBasedAgent) PlayerID() string {
	return a.playerID
}

// TODO: This is just for now, but I like being able to walk through the
// rules agent.
func (a *RuleBasedAgent) ReportState(state *game.GameState) {
}

func NewRuleBasedAgent(strategyFile string, playerID string) (*RuleBasedAgent, error) {
	agent := RuleBasedAgent{
		// TODO : make this configurable
		verbose:  true,
		playerID: playerID,
	}
	strategy, err := LoadStrategy(strategyFile)
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

func (a *RuleBasedAgent) GetNextAction(state *game.GameState) *game.GameAction {
	player, err := state.GetPlayer(a.playerID)
	if err != nil {
		// TODO: handle this more gracefully
		panic("player not found in game state")
	}
	ctx := &EvaluatorContext{
		state:    state,
		player:   player,
		strategy: a.strategy,
	}
	for _, mode := range a.strategy.Modes {
		fmt.Println("Evaluating mode:", mode.Name)
		if mode.Name == player.Mode {
			continue // Skip the current mode
		}
		result, err := mode.When.Evaluate(ctx)
		if err != nil {
			fmt.Println("Error evaluating condition for mode:", mode.Name, "-", err)
			panic("failed to evaluate condition for mode")
			return nil
		}
		if result {
			fmt.Println("Matched mode:", mode.Name)
			player.Mode = mode.Name
			fmt.Println("We did it! Switched to mode:", player.Mode)
			break
		}
		fmt.Println("Mode", mode.Name, "did not match")
	}
	var action *game.GameAction
	for _, rule := range a.strategy.Rules[player.Mode] {
		fmt.Println("Evaluating rule:", rule.Name)
		result, err := rule.When.Evaluate(ctx)
		if err != nil {
			fmt.Println("Error evaluating condition for rule:", rule.Name, "-", err)
			panic("failed to evaluate condition for rule")
			return nil
		}
		if result {
			fmt.Println("!!!!!!! Matched rule:", rule.Name)
			fmt.Println("then =>", rule.Then)
			action, err = rule.Then.Resolve(ctx)
			if err != nil {
				fmt.Println("Error evaluating action for rule:", rule.Name, "-", err)
				a.debugPrintStuff(state)
				return &game.GameAction{Type: game.ActionPass} // If there's an error, just pass}
			}
			break // Stop after the first matching rule
		}
	}
	if action == nil {
		fmt.Println("@@@@@@ No action matched in mode:", player.Mode)
		a.debugPrintStuff(state)
		return &game.GameAction{Type: game.ActionPass} // No action matched, just pass
	}
	a.debugPrintStuff(state)
	return action
}

func (a *RuleBasedAgent) EnterNumber(string, game.ChoiceSource) (int, error) {
	// For now, just return a fixed number
	return 0, nil
}
