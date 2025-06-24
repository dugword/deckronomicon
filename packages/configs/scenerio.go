package configs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// TODO: Dynamically check for yaml or json

type DeckList struct {
	Name  string `json:"Name" yaml:"Name"`
	Cards []struct {
		Name  string `json:"Name" yaml:"Name"`
		Count int    `json:"Count" yaml:"Count"`
	} `json:"Cards" yaml:"Cards"`
}

type Scenario struct {
	CheatsEnabled bool
	MaxTurns      int
	OnThePlay     string
	Name          string
	Setup         Setup
	Players       map[string]Player
}

type Player struct {
	AgentType    string
	DeckList     DeckList
	Name         string
	StartingHand []string
	StartingLife int
	StartingMode string
	StrategyFile string
}

type Setup struct {
	CheatsEnabled bool   `json:"CheatsEnabled" yaml:"CheatsEnabled"`
	MaxTurns      int    `json:"MaxTurns" yaml:"MaxTurns"`
	OnThePlay     string `json:"OnThePlay" yaml:"OnThePlay"`
	Players       []struct {
		Agent        string   `json:"Agent" yaml:"Agent"`
		DeckList     string   `json:"DeckList" yaml:"DeckList"`
		Name         string   `json:"Name" yaml:"Name"`
		StartingHand []string `json:"StartingHand" yaml:"StartingHand"`
		StartingLife int      `json:"StartingLife" yaml:"StartingLife"`
		StartingMode string   `json:"StartingMode" yaml:"StartingMode"`
		Strategy     string   `json:"Strategy" yaml:"Strategy"`
	} `json:"Players" yaml:"Players"`
	ScenarioName string `json:"ScenarioName" yaml:"ScenarioName"`
}

func LoadScenario(scenariosDir, scenario string) (*Scenario, error) {
	scenarioPath := path.Join(scenariosDir, scenario)
	setupFile := path.Join(scenarioPath, "setup.yaml")
	setupData, err := os.ReadFile(setupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read setup file %q: %w", setupFile, err)
	}
	var setup Setup
	if err := yaml.Unmarshal(setupData, &setup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setup data %q: %w", setupFile, err)
	}
	if len(setup.Players) == 0 {
		return nil, fmt.Errorf("scenario requires at least one player")
	}
	s := Scenario{
		CheatsEnabled: setup.CheatsEnabled,
		MaxTurns:      setup.MaxTurns,
		OnThePlay:     setup.OnThePlay,
		Setup:         setup,
		Name:          scenario,
		Players:       map[string]Player{},
	}
	if setup.ScenarioName != "" {
		s.Name = s.Setup.ScenarioName
	}
	// Set Default values if not provided
	if s.MaxTurns == 0 {
		s.MaxTurns = 100
	}
	var isOnThePlaySet = false
	for i, playerSetup := range setup.Players {
		if playerSetup.DeckList == "" {
			return nil, fmt.Errorf(
				"player %s requires a decklist",
				playerSetup.Name,
			)
		}
		if playerSetup.Agent == "Auto" && playerSetup.Strategy == "" {
			return nil, fmt.Errorf(
				"player %s requires a strategy file for Auto agent",
				playerSetup.Name,
			)
		}
		if setup.OnThePlay == setup.Players[i].Name {
			isOnThePlaySet = true
		}
		var player = Player{
			AgentType:    playerSetup.Agent,
			Name:         playerSetup.Name,
			StartingHand: playerSetup.StartingHand,
			StartingLife: playerSetup.StartingLife,
			StartingMode: playerSetup.StartingMode,
			StrategyFile: playerSetup.Strategy,
		}
		if player.AgentType == "" {
			player.AgentType = "Dummy"
		}
		if player.Name == "" {
			player.Name = fmt.Sprintf("Player %d", i+1)
		}
		if player.StartingLife == 0 {
			player.StartingLife = 20
		}
		if player.StartingMode == "" {
			player.StartingMode = "Setup"
		}
		if playerSetup.Strategy != "" {
			fileName := playerSetup.Strategy
			if !strings.HasSuffix(playerSetup.Strategy, ".yaml") {
				fileName += ".yaml"
			}
			player.StrategyFile = path.Join(scenarioPath, fileName)
		}
		deckList, err := LoadDeckList(scenarioPath, playerSetup.DeckList)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to load decklist for player %q: %w",
				playerSetup.Name,
				err,
			)
		}
		player.DeckList = *deckList
		if _, ok := s.Players[player.Name]; ok {
			return nil, fmt.Errorf(
				"player name %q is not unique in scenario %q",
				player.Name,
				scenario,
			)
		}
		s.Players[player.Name] = player
	}
	if !isOnThePlaySet {
		setup.OnThePlay = setup.Players[0].Name
	}
	return &s, nil
}

func LoadDeckList(scenarioPath, deckListFile string) (*DeckList, error) {
	var deckList DeckList
	fileName := deckListFile
	if !strings.HasSuffix(deckListFile, ".yaml") {
		fileName += ".yaml"
	}
	deckListPath := path.Join(scenarioPath, fileName)
	deckListData, err := os.ReadFile(deckListPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read decklist %q: %w",
			deckListPath,
			err,
		)
	}
	if err := yaml.Unmarshal(deckListData, &deckList); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal decklist %q: %w",
			deckListPath,
			err,
		)
	}
	if deckList.Name == "" {
		deckList.Name = "Unnammed Deck"
	}
	if len(deckList.Cards) == 0 {
		return nil, errors.New("decklist empty")
	}
	return &deckList, nil
}
