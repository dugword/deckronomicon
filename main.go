package main

import (
	"bufio"
	"context"
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/agent/interactive"
	"deckronomicon/packages/configs"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/logger"
	"deckronomicon/packages/state"
	"fmt"
	"io"
	"os"
)

// TODO: Remove this seed and make it configurable.
const seed = 13

func createPlayerAgent(
	playerScenario configs.Player,
	config configs.Config,
	stdin io.Reader,
) (engine.PlayerAgent, error) {
	var playerAgent engine.PlayerAgent
	switch playerScenario.AgentType {
	case "Interactive":
		scanner := bufio.NewScanner(stdin)
		playerAgent = interactive.NewAgent(
			scanner,
			playerScenario.Name,
			[]mtg.Step{mtg.StepPrecombatMain},
			config.Verbose,
		)
		return playerAgent, nil
	//case "Auto":
	/*
		var err error
		playerAgent, err = auto.NewRuleBasedAgent(
			playerScenario.StrategyFile,
			config.Interactive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule based agent: %w", err)
		}
	*/
	case "Dummy":
		playerAgent = dummy.NewAgent(
			playerScenario.Name,
			[]mtg.Step{mtg.StepPrecombatMain},
			config.Verbose,
		)
		return playerAgent, nil
	default:
		return nil, fmt.Errorf(
			"unknown player agent type '%s'",
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
	// Logger for application level information.
	logger := logger.NewLogger()
	logger.Info("Starting Deckronomicon...")
	logger.Info("Loading scenario...")
	scenario, err := configs.LoadScenario(config.ScenariosDir, config.Scenario)
	if err != nil {
		return fmt.Errorf("failed to load scenario: %w", err)
	}
	if config.Cheat {
		scenario.Setup.CheatsEnabled = true
	}
	logger.Info(fmt.Sprintf("Scenario '%s' loaded!", scenario.Name))
	logger.Info("Loading card definitions...")
	cardDefinitions, err := definition.LoadCardDefinitions(config.Definitions)
	if err != nil {
		return fmt.Errorf("failed to load card definitions: %w", err)
	}
	logger.Info("Card definitions loaded!")
	playerAgents := map[string]engine.PlayerAgent{}
	for _, playerScenario := range scenario.Players {
		logger.Info(fmt.Sprintf(
			"Creating player '%s' with agent type '%s'...",
			playerScenario.Name,
			playerScenario.AgentType,
		))
		playerAgent, err := createPlayerAgent(playerScenario, config, stdin)
		if err != nil {
			return fmt.Errorf("failed to create player agent for '%s': %w", playerScenario.Name, err)
		}
		playerAgents[playerScenario.Name] = playerAgent
		logger.Info(fmt.Sprintf("Player Agent for '%s' created!", playerScenario.Name))
	}
	deckLists := map[string]configs.DeckList{}
	for _, playerScenario := range scenario.Players {
		deckLists[playerScenario.Name] = playerScenario.DeckList
	}
	var players []state.Player
	for _, playerScenario := range scenario.Players {
		logger.Info(fmt.Sprintf(
			"Creating player '%s' with starting life %d...",
			playerScenario.Name,
			playerScenario.StartingLife,
		))
		player := state.NewPlayer(playerScenario.Name, playerScenario.StartingLife)
		players = append(players, player)
		logger.Info(fmt.Sprintf("Player '%s' created!", player.ID()))
	}
	engineConfig := engine.EngineConfig{
		Agents:      playerAgents,
		Definitions: cardDefinitions,
		DeckLists:   deckLists,
		Players:     players,
		Seed:        seed,
	}
	engine := engine.NewEngine(engineConfig)
	if err := engine.RunGame(); err != nil {
		return fmt.Errorf("failed to run the game: %w", err)
	}
	logger.Info("Game completed successfully!")
	return nil
}
