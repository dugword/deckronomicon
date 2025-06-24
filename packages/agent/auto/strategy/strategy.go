package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/state"
)

type Rule struct {
	Name        string              `json:"Name" yaml:"Name"`
	Description string              `json:"Description" yaml:"Description"`
	RawWhen     map[string]any      `json:"When" yaml:"When"`
	When        evaluator.Evaluator `json:"-" yaml:"-"`
	RawThen     map[string]any      `json:"Then" yaml:"Then"`
	Then        action.ActionNode   `json:"-" yaml:"-"`
}

type Strategy struct {
	Name        string              `json:"Name,omitempty" yaml:"Name,omitempty"`
	Description string              `json:"Description,omitempty" yaml:"Description,omitempty"`
	Definitions map[string][]string `json:"Definitions,omitempty" yaml:"Definitions,omitempty"`
	Modes       []Rule              `json:"Modes,omitempty" yaml:"Modes,omitempty"`
	Rules       map[string][]Rule   `json:"Rules,omitempty" yaml:"Rules,omitempty"`
}

type EvaluatorContext struct {
	Game     state.Game
	Player   state.Player
	Strategy *Strategy
}
