package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Setup struct {
	MaxTurns             int      `json:"MaxTurns"`
	OnThePlay            bool     `json:"OnThePlay"`
	OpponentName         string   `json:"OpponentName"`
	OpponentStartingHand []string `json:"OpponentStartingHand"`
	PlayerName           string   `json:"PlayerName"`
	PlayerStartingHand   []string `json:"StartingHand"`
	StartingLife         int      `json:"StartingLife"`
}

type Scenerio struct {
	Name             string
	OpponentDeck     string
	OpponentStrategy string
	PlayerDeck       string
	PlayerStrategy   string
	Setup            Setup
}

// DeckImportCard represents a card in the deck import file.
type DeckImportCard struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// DeckImport represents the structure of the deck import file.
type DeckImport struct {
	Cards []DeckImportCard `json:"cards"`
}

func LoadScenario(scenariosDir, scenario string) (*Scenerio, error) {
	scenariosPath := path.Join(scenariosDir, scenario)
	s := Scenerio{
		Name:             scenario,
		OpponentDeck:     path.Join(scenariosPath, "opponent_deck.json"),
		OpponentStrategy: path.Join(scenariosPath, "opponent_strategy.json"),
		PlayerDeck:       path.Join(scenariosPath, "player_deck.json"),
		PlayerStrategy:   path.Join(scenariosPath, "player_strategy.json"),
	}
	setupData, err := os.ReadFile(path.Join(scenariosPath, "setup.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read setup file: %w", err)
	}
	if err := json.Unmarshal(setupData, &s.Setup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setup data: %w", err)
	}
	return &s, nil
}
