package auto

import (
	"deckronomicon/packages/agent/auto/strategy"
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

func (a *RuleBasedAgent) ReportState(game *state.Game) {
	player := game.GetPlayer(a.playerID)
	opponent := game.GetOpponent(a.playerID)
	a.uiBuffer.Update(
		view.NewGameViewFromState(game),
		view.NewPlayerViewFromState(game, player.ID(), a.mode),
		view.NewPlayerViewFromState(game, opponent.ID(), ""),
	)
}

func NewRuleBasedAgent(
	scenarioPath string,
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
	strategy, err := strategy.LoadStrategy(
		scenarioPath,
		strategyFile,
		strategy.StrategyDirectories{
			Actions:    "actions",
			Choices:    "choices",
			Conditions: "conditions",
			Groups:     "groups",
			Modes:      "modes",
			Rules:      "rules",
			Selectors:  "selectors",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy: %w", err)
	}
	agent.strategy = strategy
	return &agent, nil
}

func (a *RuleBasedAgent) PlayerID() string {
	return a.playerID
}
