package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"deckronomicon/auto"
	"deckronomicon/game"
	"deckronomicon/interactive"
	"deckronomicon/logger"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	flags := flag.NewFlagSet("deckronomicon", flag.ContinueOnError)
	isInteractive := flags.Bool("interactive", false, "run as interactive mode")
	maxTurns := flags.Int("max-turns", 100, "maximum number of turns to simulate")
	err := flags.Parse(args[1:])
	if err != nil {
		// TODO: Handle this without fmt
		return err
	}

	// TODO: Handle errors as a separate writer
	logger := logger.NewLogger(w, w)
	logger.Log("Starting Deckronomicon...")

	logger.Log("Creating new game state...")
	state := game.NewGameState()
	logger.Log("New game state created!")

	logger.Log("Initializing new game...")
	if err := state.InitializeNewGame(game.GameStateConfig{
		DeckList:     "my_deck.json",
		MaxTurns:     *maxTurns,
		StartingLife: 20,
	}); err != nil {
		return err
	}
	logger.Log("New game initalized!")

	var playerAgent game.PlayerAgent
	if *isInteractive {
		logger.Log("Running in interactive mode")
		// TODO: Pass this into main so we can test?
		// maybe not because we can test the interative package separately...
		scanner := bufio.NewScanner(os.Stdin)
		playerAgent = interactive.NewInteractivePlayerAgent(scanner)
	} else {
		logger.Log("Running an auto simulation")
		playerAgent = auto.NewRuleBasedAgent("sample_rules.json")
		//playerAgent = auto.NewAutoPlayerAgent()
	}

	gameEndErr := state.RunGameLoop(playerAgent)
	for _, message := range state.MessageLog {
		fmt.Println(message)
	}
	if gameEndErr != nil {
		return gameEndErr
	}
	return nil
}
