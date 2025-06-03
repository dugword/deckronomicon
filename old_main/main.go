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
	stdin io.Reader,
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
	logger.Log("Creating new players...")
	var players []*player.Player
	for _, playerScenario := range scenario.Players {
		logger.Log(fmt.Sprintf(
			"Creating player '%s' with agent type '%s'...\n",
			playerScenario.Name,
			playerScenario.AgentType,
		))
		p, err := createPlayer(playerScenario, config, stdin)
		if err != nil {
			return fmt.Errorf("failed to create player: %s", playerScenario.Name)
		}
		players = append(players, p)
		logger.Log(fmt.Sprintf("Player '%s' created!", p.ID()))
	}
	logger.Log("All players created!")
	logger.Log("Initializing new game...")
	engine := engine.NewGame(
		scenario,
		players,
		cardDefinitions,
	)
	logger.Log("New game initialized!")
	logger.Log("Running game loop...")
	err = engine.RunGameLoop()
	// TODO: Split game state object from engine object.
	err = engine.RunGameLoop()
	if err != nil && !errors.Is(err, mtg.ErrGameOver) {
		return fmt.Errorf("game loop failed: %w", err)
	}
	// TODO: This sucks
	logger.Log("Game Message Log:\n" + strings.Join(engine.GameState.MessageLog, "\n"))
	logger.Log("Game over!")
	return nil
}

func createPlayer(
	playerScenario configs.Player,
	config configs.Config,
	stdin io.Reader,
) (*player.Player, error) {
	var playerAgent player.Agent
	switch playerScenario.AgentType {
	case "Interactive":
		scanner := bufio.NewScanner(stdin)
		playerAgent = interactive.NewInteractivePlayerAgent(
			scanner,
		)
	case "Auto":
		var err error
		playerAgent, err = auto.NewRuleBasedAgent(
			playerScenario.StrategyFile,
			config.Interactive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule based agent: %w", err)
		}
	case "Dummy":
		playerAgent = dummy.NewDummyAgent()
	default:
		return nil, fmt.Errorf(
			"unknown player agent type: %s",
			playerScenario.AgentType,
		)
	}
	p := player.New(
		playerAgent,
		playerScenario.Name,
		playerScenario.StartingLife,
		playerScenario.StartingMode,
	)
	return p, nil
}
