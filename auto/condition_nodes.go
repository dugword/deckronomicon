package auto

import (
	"deckronomicon/game"
	"fmt"
)

type ConditionNode interface {
	Evaluate(ctx *EvaluatorContext) (bool, error)
}

// Card condition
type CardConditionNode struct {
	Zone  string        `json:"zone"`
	Cards CardCondition `json:"cards"`
}

type CardNameCondition struct {
	Name string `json:"name"`
}

func (c *CardNameCondition) Matches(objs []game.GameObject, definitions map[string][]string) bool {
	for _, obj := range objs {
		if obj.Name() == c.Name {
			return true
		}
	}
	return false
}

func (n *CardConditionNode) Evaluate(ctx *EvaluatorContext) (bool, error) {
	return true, nil
}

// AND condition
type AndCondition struct {
	Conditions []ConditionNode
}

func (c *AndCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	for _, cond := range c.Conditions {
		result, err := cond.Evaluate(ctx)
		if err != nil {
			return false, fmt.Errorf("error evaluating condition: %w", err)
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

// OR condition
type OrCondition struct {
	Conditions []ConditionNode
}

func (c *OrCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	for _, cond := range c.Conditions {
		result, err := cond.Evaluate(ctx)
		if err != nil {
			return false, fmt.Errorf("error evaluating condition: %w", err)
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

// NOT condition
type NotCondition struct {
	Condition ConditionNode
}

func (c *NotCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	result, err := c.Condition.Evaluate(ctx)
	if err != nil {
		return false, fmt.Errorf("error evaluating condition: %w", err)
	}
	return !result, nil
}

type StepCondition struct {
	Step string `json:"Step"`
}

func (c *StepCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	// You’ll need to implement actual step checking logic here
	// For now, just return true if the step matches
	if ctx.state.CurrentStep == c.Step {
		return true, nil
	}
	return false, nil
}

type ModeCondition struct {
	Mode string `json:"Mode"`
}

func (c *ModeCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	// You’ll need to implement actual mode checking logic here
	// For now, just return true if the mode matches
	if ctx.player.Mode == c.Mode {
		return true, nil
	}
	return false, nil
}

// InZone condition
type InZoneCondition struct {
	Zone  string `json:"Zone"`
	Cards CardCondition
}

func (c *InZoneCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	if c.Cards == nil {
		return true, nil // No specific cards means always true
	}
	zone, err := ctx.player.GetZone(c.Zone)
	if err != nil {
		return false, fmt.Errorf("error getting zone '%s': %w", c.Zone, err)
	}
	if c.Cards.Matches(zone.GetAll(), ctx.strategy.Definitions) {
		return true, nil // Zone contains the specified cards
	}
	return false, nil
}

// PlayerStat condition
type PlayerStatCondition struct {
	Stat  string `json:"Stat"`
	Op    string `json:"Op"`
	Value int    `json:"Value"`
}

func (c *PlayerStatCondition) Evaluate(ctx *EvaluatorContext) (bool, error) {
	// You’ll need to implement actual stat checking logic here
	return false, nil
}
