package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"deckronomicon/auto"
	"deckronomicon/configs"
	"deckronomicon/dummy"
	"deckronomicon/game"
	"deckronomicon/interactive"
	"deckronomicon/log"
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
	logger.Log("Scenario loaded!")
	logger.Log("Creating player agent...")
	var playerAgent game.PlayerAgent
	if config.Interactive {
		logger.Log("Creating interactive player agent...")
		scanner := bufio.NewScanner(stdin)
		playerAgent = interactive.NewInteractivePlayerAgent(
			scanner,
			scenario.Setup.PlayerName,
		)
	} else {
		logger.Log("Creating rule based player agent...")
		playerAgent = auto.NewRuleBasedAgent(
			scenario.PlayerStrategy,
			scenario.Setup.PlayerName,
		)
	}
	opponentAgent := dummy.NewDummyAgent(
		scenario.Setup.OpponentName,
	)
	logger.Log("Player agent created!")
	logger.Log("Creating new players...")
	player := game.NewPlayer(
		playerAgent,
	)
	// TODO I don't like this, but it works for now
	player.StartingHand = scenario.Setup.PlayerStartingHand
	opponent := game.NewPlayer(
		opponentAgent,
	)
	// TODO I don't like this, but it works for now
	opponent.StartingHand = scenario.Setup.OpponentStartingHand
	logger.Log("Players created!")
	logger.Log("Creating new game state...")
	state := game.NewGameState()
	logger.Log("New game state created!")
	logger.Log("New players created!")
	logger.Log("Initializing new game...")
	if err := state.InitializeNewGame(
		scenario,
		player,
		opponent,
		config,
	); err != nil {
		return fmt.Errorf("failed to initialize new game: %w", err)
	}
	logger.Log("New game initalized!")
	logger.Log("Running game loop...")
	// TOOD: This still tracks game losses with application errors
	// Those should be separated out
	err = state.RunGameLoop()
	logger.Log("Game Message Log:\n" + strings.Join(state.MessageLog, "\n"))
	if err != nil {
		return fmt.Errorf("game loop failed: %w", err)
	}
	return nil
}
