package evaluator

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/game/mtg"
)

type Evaluator interface {
	Evaluate(ctx *evalstate.EvalState) bool
}

type True struct{}

func (e *True) Evaluate(ctx *evalstate.EvalState) bool {
	return true
}

type False struct{}

func (e *False) Evaluate(ctx *evalstate.EvalState) bool {
	return false
}

type And struct {
	Evaluators []Evaluator
}

func (e *And) Evaluate(ctx *evalstate.EvalState) bool {
	for _, evaluator := range e.Evaluators {
		if !evaluator.Evaluate(ctx) {
			return false
		}
	}
	return true
}

type Or struct {
	Evaluators []Evaluator
}

func (e *Or) Evaluate(ctx *evalstate.EvalState) bool {
	for _, evaluator := range e.Evaluators {
		if evaluator.Evaluate(ctx) {
			return true
		}
	}
	return false
}

type Not struct {
	Evaluator Evaluator
}

func (e *Not) Evaluate(ctx *evalstate.EvalState) bool {
	return !e.Evaluator.Evaluate(ctx)
}

type Step struct {
	Step mtg.Step
}

func (e *Step) Evaluate(ctx *evalstate.EvalState) bool {
	return ctx.Game.Step() == e.Step
}

type LandPlayedThisTurn struct {
	LandPlayedThisTurn bool
}

func (e *LandPlayedThisTurn) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	return player.LandPlayedThisTurn() == e.LandPlayedThisTurn
}

type Mode struct {
	Mode string
}

func (e *Mode) Evaluate(ctx *evalstate.EvalState) bool {
	return ctx.Mode == e.Mode
}

type PlayerStat struct {
	Stat  string
	Op    string
	Value int
}

func (e *PlayerStat) Evaluate(ctx *evalstate.EvalState) bool {
	return false
}
