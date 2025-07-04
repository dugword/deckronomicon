package strategy

import "testing"

func TestLoadStrategy(t *testing.T) {
	strategy, err := LoadStrategy(
		"testdata",
		"player_strategy.yaml",
		StrategyDirectories{
			Actions:    "actions",
			Choices:    "choices",
			Conditions: "conditions",
			Groups:     "groups",
			Modes:      "modes",
			Rules:      "rules",
			Selectors:  "selectors",
		},
	)
	if err != nil {
		t.Fatalf("Failed to load strategy: %v", err)
	}

	if strategy.Name != "Test Player Strategy" {
		t.Errorf("Expected strategy name 'Player Strategy', got '%s'", strategy.Name)
	}

	if len(strategy.Rules) == 0 {
		t.Error("Expected at least one rule in the strategy")
	}

	if len(strategy.Choices) == 0 {
		t.Error("Expected at least one choice in the strategy")
	}
}
