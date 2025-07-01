package main

import (
	"bufio"
	"context"
	"deckronomicon/packages/agent/auto"
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/agent/interactive"
	"deckronomicon/packages/analytics"
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/logger"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func createPlayerAgent(
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
	logger.Info("Starting Deckronomicon...")
	logger.Info("Loading scenario...")
	scenario, err := configs.LoadScenario(config.ScenariosDir, config.Scenario)
	if err != nil {
		return fmt.Errorf("failed to load scenario: %w", err)
	}
	if config.Cheat {
		scenario.Setup.CheatsEnabled = true
	}
	logger.Info(fmt.Sprintf("Scenario %q loaded!", scenario.Name))
	logger.Info("Loading card definitions...")
	cardDefinitions, err := definition.LoadCardDefinitions(config.Definitions)
	if err != nil {
		return fmt.Errorf("failed to load card definitions: %w", err)
	}
	logger.Info("Card definitions loaded!")
	playerAgents := map[string]engine.PlayerAgent{}
	for _, playerScenario := range scenario.Players {
		logger.Info(fmt.Sprintf(
			"Creating player %q with agent type %q...",
			playerScenario.Name,
			playerScenario.AgentType,
		))
		playerAgent, err := createPlayerAgent(playerScenario, config, stdin)
		if err != nil {
			return fmt.Errorf("failed to create player agent for %q: %w", playerScenario.Name, err)
		}
		playerAgents[playerScenario.Name] = playerAgent
		logger.Info(fmt.Sprintf("Player Agent for %q created!", playerScenario.Name))
	}
	deckLists := map[string]configs.DeckList{}
	for _, playerScenario := range scenario.Players {
		deckLists[playerScenario.Name] = playerScenario.DeckList
	}
	var players []state.Player
	for _, playerScenario := range scenario.Players {
		logger.Info(fmt.Sprintf(
			"Creating player %q with starting life %d...",
			playerScenario.Name,
			playerScenario.StartingLife,
		))
		player := state.NewPlayer(playerScenario.Name, playerScenario.StartingLife)
		players = append(players, player)
		logger.Info(fmt.Sprintf("Player %q created!", player.ID()))
	}
	seed := time.Now().UnixNano()
	if scenario.Setup.Seed != 0 && config.Seed == 0 {
		logger.Info(fmt.Sprintf("Using seed %d from scenario setup.", scenario.Setup.Seed))
		seed = scenario.Setup.Seed
	}
	if config.Seed != 0 {
		logger.Info(fmt.Sprintf("Using seed %d from command line config.", config.Seed))
		seed = config.Seed
	}
	engineConfig := engine.EngineConfig{
		Agents:      playerAgents,
		Definitions: cardDefinitions,
		DeckLists:   deckLists,
		Players:     players,
		Seed:        seed,
		Log:         logger,
	}
	engine := engine.NewEngine(engineConfig)
	if err := engine.RunGame(); err != nil && !errors.Is(err, mtg.ErrGameOver) {
		return fmt.Errorf("failed to run the game: %w", err)
	}
	game := engine.Game()
	logger.Info("Game completed successfully!")
	logger.Info("Game over!")
	var playerLostErr mtg.PlayerLostError
	if errors.As(err, &playerLostErr) {
		logger.Info(fmt.Sprintf("Game over reason: %s", playerLostErr.Reason))
	}
	if err := analytics.WriteAnalyticsEventsToFile(engine.Record(), "results"); err != nil {
		return fmt.Errorf("failed to write game record: %w", err)
	}
	if err := analytics.WriteGameRecordToFile(engine.Record(), "results"); err != nil {
		return fmt.Errorf("failed to write game record: %w", err)
	}
	if err := analytics.WriteGameStateToFile(game, "results"); err != nil {
		return fmt.Errorf("failed to write game state: %w", err)
	}
	logger.Info("Game record and state written to results directory.")

	return nil
}
