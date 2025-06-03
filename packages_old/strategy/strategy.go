package strategy

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/player"
	"deckronomicon/packages/strategy/action"
	"deckronomicon/packages/strategy/condition"
)

type Rule struct {
	Name        string                  `json:"Name"`
	Description string                  `json:"Description"`
	RawWhen     map[string]any          `json:"When"`
	When        condition.ConditionNode `json:"-"`
	RawThen     map[string]any          `json:"Then"`
	Then        action.ActionNode       `json:"-"`
}

type Strategy struct {
	Name        string              `json:"Name,omitempty"`
	Description string              `json:"Description,omitempty"`
	Definitions map[string][]string `json:"Definitions,omitempty"`
	Modes       []Rule              `json:"Modes,omitempty"`
	Rules       map[string][]Rule   `json:"Rules,omitempty"`
}

type EvaluatorContext struct {
	State    *engine.GameState
	Player   *player.Player
	Strategy *Strategy
}
