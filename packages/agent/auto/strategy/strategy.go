package strategy

import (
	"deckronomicon/packages/agent/auto/strategy/action"
	"deckronomicon/packages/agent/auto/strategy/evaluator"
	"deckronomicon/packages/state"
)

type Rule struct {
	Name        string
	Description string
	When        evaluator.Evaluator
	Then        action.ActionNode
}

type Strategy struct {
	Name        string
	Description string
	Modes       []*Rule
	Rules       map[string][]*Rule
}

type EvaluatorContext struct {
	Game     *state.Game
	PlayerID string
	Strategy *Strategy
}
