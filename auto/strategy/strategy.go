package strategy

import (
	"deckronomicon/game"
)

type Rule struct {
	Name        string         `json:"Name"`
	Description string         `json:"Description"`
	RawWhen     map[string]any `json:"When"`
	When        ConditionNode  `json:"-"`
	RawThen     map[string]any `json:"Then"`
	Then        ActionNode     `json:"-"`
}

type Strategy struct {
	Name        string              `json:"Name,omitempty"`
	Description string              `json:"Description,omitempty"`
	Definitions map[string][]string `json:"Definitions,omitempty"`
	Modes       []Rule              `json:"Modes,omitempty"`
	Rules       map[string][]Rule   `json:"Rules,omitempty"`
}

type EvaluatorContext struct {
	State    *game.GameState
	Player   *game.Player
	Strategy *Strategy
}
