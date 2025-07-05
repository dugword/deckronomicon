package main

import (
	"bufio"
	"context"
	"deckronomicon/packages/agent/auto"
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/agent/interactive"
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/store"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func createPlayerAgent(
	scenario *configs.Scenario,
	playerScenario configs.Player,
	config configs.Config,
	stdin io.Reader,
) (engine.PlayerAgent, error) {
	var playerAgent engine.PlayerAgent
	switch playerScenario.AgentType {
	case "Interactive":
		scanner := bufio.NewScanner(stdin)
		autoPay := playerScenario.AutoPay
		if config.AutoPay {
			autoPay = config.AutoPay
		}
		playerAgent = interactive.NewAgent(
			scanner,
			playerScenario.Name,
			[]mtg.Step{mtg.StepUpkeep, mtg.StepPrecombatMain},
			"./ui/term/display.tmpl", // TODO: Make this configurable.
			autoPay,
			playerScenario.AutoPayColorsForGeneric,
			config.Verbose,
		)
		return playerAgent, nil
	case "Auto":
		var err error
		playerAgent, err = auto.NewRuleBasedAgent(
			scenario.DirName,
			playerScenario.StrategyFile,
			playerScenario.Name,
			"./ui/term/display.tmpl", // TODO: Make this configurable.
			config.Interactive,
			[]mtg.Step{mtg.StepPrecombatMain},
			config.Verbose,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule based agent: %w", err)
		}
		return playerAgent, nil
	case "Dummy":
		playerAgent = dummy.NewChooseMinimumAgent(
			playerScenario.Name,
		)
		return playerAgent, nil
	default:
		return nil, fmt.Errorf(
			"unknown player agent type %q",
			playerScenario.AgentType,
		)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	if err := Run(
		ctx,
		cancel,
		os.Args,
		os.Getenv,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Run is an abstraction for the main function to enable testing.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	defer cancel()
	config, err := configs.LoadConfig(args, getenv)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger := logger.NewLogger()
	logger.LogLevel = config.LogLevel
	logger.Debug("Starting Deckronomicon...")
	logger.Debug("Loading scenario...")
	scenario, err := configs.LoadScenario(config.ScenariosDir, config.Scenario)
	if err != nil {
		return fmt.Errorf("failed to load scenario: %w", err)
	}
	if config.Cheat {
		scenario.Setup.CheatsEnabled = true
	}
	logger.Debugf("Scenario %q loaded!", scenario.Name)
	logger.Debug("Loading card definitions...")
	cardDefinitions, err := definition.LoadCardDefinitions(config.Definitions)
	if err != nil {
		return fmt.Errorf("failed to load card definitions: %w", err)
	}
	logger.Debug("Card definitions loaded!")
	playerAgents := map[string]engine.PlayerAgent{}
	for _, playerScenario := range scenario.Players {
		logger.Debug(fmt.Sprintf(
			"Creating player %q with agent type %q...",
			playerScenario.Name,
			playerScenario.AgentType,
		))
		playerAgent, err := createPlayerAgent(scenario, playerScenario, config, stdin)
		if err != nil {
			return fmt.Errorf("failed to create player agent for %q: %w", playerScenario.Name, err)
		}
		playerAgents[playerScenario.Name] = playerAgent
		logger.Debugf("Player Agent for %q created!", playerScenario.Name)
	}
	deckLists := map[string]*configs.DeckList{}
	for _, playerScenario := range scenario.Players {
		deckLists[playerScenario.Name] = playerScenario.DeckList
	}
	var playerDefinitions []*definition.Player
	for _, playerScenario := range scenario.Players {
		logger.Debugf("Creating player %q with starting life %d...",
			playerScenario.Name,
			playerScenario.StartingLife,
		)
		playerDefinition := &definition.Player{
			ID:          playerScenario.Name,
			Life:        playerScenario.StartingLife,
			MaxHandSize: 7, // Default max hand size
		}
		playerDefinitions = append(playerDefinitions, playerDefinition)
		logger.Debugf("Player definition %q created!", playerDefinition.ID)
	}
	seed := time.Now().UnixNano()
	if scenario.Setup.Seed != 0 && config.Seed == 0 {
		logger.Debugf("Using seed %d from scenario setup.", scenario.Setup.Seed)
		seed = scenario.Setup.Seed
	}
	if config.Seed != 0 {
		logger.Debugf("Using seed %d from command line config.", config.Seed)
		seed = config.Seed
	}
	engineConfig := engine.EngineConfig{
		Agents:            playerAgents,
		CardDefinitions:   cardDefinitions,
		DeckLists:         deckLists,
		PlayerDefinitions: playerDefinitions,
		Seed:              seed,
		Log:               logger,
	}
	start := time.Now()

	runResults := []store.RunResult{}
	for i := range config.Runs {
		runID := i + 1
		runResult := store.RunResult{
			RunID:       runID,
			Totals:      map[string]int{},
			Cumulatives: map[string][]int{},
		}
		engine := engine.NewEngine(runID, &runResult, engineConfig)
		if err := engine.RunGame(); err != nil && !errors.Is(err, mtg.ErrGameOver) {
			return fmt.Errorf("failed to run the game: %w", err)
		}
		engineConfig.Seed++
		logger.Debug("Game over!")
		logger.Infof("Game %d completed successfully!", i+1)
		runResults = append(runResults, runResult)
	}
	end := time.Now()
	logger.Infof("Total time taken for %d runs: %s", config.Runs, end.Sub(start))
	logger.Infof("Average time per run: %s", end.Sub(start)/time.Duration(config.Runs))
	out, err := json.MarshalIndent(runResults, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	if err := os.WriteFile("results/results.json", out, 0644); err != nil {
		return fmt.Errorf("failed to write results to file: %w", err)
	}
	return nil
}
