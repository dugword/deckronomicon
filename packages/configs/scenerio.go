package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type DeckList struct {
	Name  string `json:"Name"`
	Cards []struct {
		Name  string `json:"Name"`
		Count int    `json:"Count"`
	} `json:"Cards"`
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
	CheatsEnabled bool   `json:"CheatsEnabled"`
	MaxTurns      int    `json:"MaxTurns"`
	OnThePlay     string `json:"OnThePlay"`
	Players       []struct {
		Agent        string   `json:"Agent"`
		DeckList     string   `json:"DeckList"`
		Name         string   `json:"Name"`
		StartingHand []string `json:"StartingHand"`
		StartingLife int      `json:"StartingLife"`
		StartingMode string   `json:"StartingMode"`
		Strategy     string   `json:"Strategy"`
	} `json:"Players"`
	ScenarioName string `json:"ScenarioName"`
}

func LoadScenario(scenariosDir, scenario string) (*Scenario, error) {
	scenarioPath := path.Join(scenariosDir, scenario)
	setupFile := path.Join(scenarioPath, "setup.json")
	setupData, err := os.ReadFile(setupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read setup file %q: %w", setupFile, err)
	}
	var setup Setup
	if err := json.Unmarshal(setupData, &setup); err != nil {
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
			if !strings.HasSuffix(playerSetup.Strategy, ".json") {
				fileName += ".json"
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
	if !strings.HasSuffix(deckListFile, ".json") {
		fileName += ".json"
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
	if err := json.Unmarshal(deckListData, &deckList); err != nil {
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
