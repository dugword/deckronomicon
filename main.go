package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"deckronomicon/packages/agent/auto"
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/agent/interactive"
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/log"
)

// main is the entry point for the application.
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

// Run is an abtraction for the main function to enable testing.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	args []string,
	getenv func(string) string,
	stdin *os.File,
	stdout io.Writer,
	stderr io.Writer,
) error {
	defer cancel()
	config, err := configs.LoadConfig(args, getenv)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	// Logger for application level information.
	logger := log.NewLogger(
		"✨",
		stdout,
		stderr,
		config.Verbose,
	)
	logger.Log("Starting Deckronomicon...")
	logger.Log("Loading scenario...")
	scenario, err := configs.LoadScenario(config.ScenariosDir, config.Scenario)
	if err != nil {
		return fmt.Errorf("failed to load scenario: %w", err)
	}
	if config.Cheat == true {
		scenario.Setup.CheatsEnabled = true
	}
	logger.Log(fmt.Sprintf("Scenario '%s' loaded!", scenario.Name))
	logger.Log("Loading card definitions...")
	cardDefinitions, err := definition.LoadCardDefinitions(config.Definitions)
	if err != nil {
		return fmt.Errorf("failed to load card definitions: %w", err)
	}
	logger.Log("Card definitions loaded!")
	var players []*player.Player
	for _, playerScenario := range scenario.Players {
		logger.Log("Creating player agents...")
		var playerAgent = player.Agent()
		switch playerScenario.AgentType {
		case "Interactive":
			logger.Log("Creating interactive player agent...")
			scanner := bufio.NewScanner(stdin)
			playerAgent = interactive.NewInteractivePlayerAgent(
				scanner,
			)
		case "Auto":
			logger.Log("Creating rule based player agent...")
			var err error
			playerAgent, err = auto.NewRuleBasedAgent(
				playerScenario.StrategyFile,
				config.Interactive,
			)
			if err != nil {
				return fmt.Errorf("failed to create rule based agent: %w", err)
			}
		case "Dummy":
			logger.Log("Creating dummy player agent...")
			playerAgent = dummy.NewDummyAgent()
		default:
			return fmt.Errorf("unknown player agent type: %s", playerSetup.Agent)
		}
		logger.Log("Player agents created!")
		logger.Log("Creating new players...")
		plyr := player.New(
			playerAgent,
			playerScenario.Name,
			playerScenario.StartingLife,
			playerScenario.StartingMode,
		)
		logger.Log("Players created!")
		players = append(players, plyr)
	}
	logger.Log("Initializing new game...")
	if err := engine.InitializeNewGame(
		scenario,
		players,
		cardDefinitions,
	); err != nil {
		return fmt.Errorf("failed to initialize new game state: %w", err)
	}
	logger.Log("New game initialized!")
	logger.Log("Running game loop...")
	err = state.RunGameLoop()
	logger.Log("Game Message Log:\n" + strings.Join(state.MessageLog, "\n"))
	// TODO: Split game state object from engine object.
	err = state.RunGameLoop()
	if err != nil && !errors.Is(err, mtg.ErrGameOver) {
		return fmt.Errorf("game loop failed: %w", err)
	}
	logger.Log("Game over!")
	return nil
}
