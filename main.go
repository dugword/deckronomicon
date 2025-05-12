package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"deckronomicon/auto"
	"deckronomicon/game"
	"deckronomicon/interactive"
	"deckronomicon/logger"
)

var INTERACTIVE = false
var MAX_TURNS = 100

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "interactive" {
			INTERACTIVE = true
		}
	}
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// TODO: Handle errors as a separate writer
	logger := logger.NewLogger(w, w)
	logger.Log("Starting Deckronomicon...")

	logger.Log("Creating new game state...")
	state := game.NewGameState()
	logger.Log("New game state created!")

	logger.Log("Initializing new game...")
	if err := state.InitializeNewGame("my_deck.json"); err != nil {
		return err
	}
	logger.Log("New game initalized!")

	var playerAgent game.PlayerAgent
	if INTERACTIVE {
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

	state.RunGameLoop(playerAgent, MAX_TURNS)
	for _, message := range state.MessageLog {
		fmt.Println(message)
	}
	return nil
}
