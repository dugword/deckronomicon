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
		"deckronomicon",
		stdout,
		stderr,
		config.Verbose,
	)
	logger.Log("Starting Deckronomicon...")
	logger.Log("Creating new game state...")
	state := game.NewGameState()
	logger.Log("New game state created!")
	logger.Log("Initializing new game...")
	if err := state.InitializeNewGame(config); err != nil {
		return fmt.Errorf("failed to initialize new game: %w", err)
	}
	logger.Log("New game initalized!")

	logger.Log("Creating player agent...")
	var playerAgent game.PlayerAgent
	if config.Interactive {
		logger.Log("Creating interactive player agent...")
		scanner := bufio.NewScanner(stdin)
		playerAgent = interactive.NewInteractivePlayerAgent(scanner)
	} else {
		logger.Log("Creating rule based player agent...")
		playerAgent = auto.NewRuleBasedAgent(config.StrategyFile)
	}
	logger.Log("Player agent created!")
	logger.Log("Running game loop...")
	// TOOD: This still tracks game losses with application errors
	// Those should be separated out
	err = state.RunGameLoop(playerAgent)
	logger.Log("Game Message Log:\n" + strings.Join(state.MessageLog, "\n"))
	if err != nil {
		return fmt.Errorf("game loop failed: %w", err)
	}
	return nil
}
