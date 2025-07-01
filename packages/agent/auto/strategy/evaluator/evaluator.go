package evaluator

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/game/mtg"
)

type Op string

const (
	OpEqual              Op = "=="
	OpGreaterThan        Op = ">"
	OpGreaterThanOrEqual Op = ">="
	OpLessThan           Op = "<"
	OpLessThanOrEqual    Op = "<="
	OpNotEqual           Op = "!="
)

func Operators() []Op {
	return []Op{
		OpEqual,
		OpGreaterThan,
		OpGreaterThanOrEqual,
		OpLessThan,
		OpLessThanOrEqual,
		OpNotEqual,
	}
}

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
	// TODO: Have this check for player turn be more explicit
	// TODO: Maybe have an "ActivePlayerStep" or
	// a flag to indicate we should also check the opponent's step
	if ctx.Game.ActivePlayerID() != ctx.PlayerID {
		return false
	}
	return ctx.Game.Step() == e.Step
}

type LandPlayedThisTurn struct {
	LandPlayedThisTurn bool
}

func (e *LandPlayedThisTurn) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	return player.LandPlayedThisTurn() == e.LandPlayedThisTurn
}

type StackEmpty struct {
	StackEmpty bool
}

func (e *StackEmpty) Evaluate(ctx *evalstate.EvalState) bool {
	return (ctx.Game.Stack().Size() == 0) == e.StackEmpty
}

type SorcerySpeed struct {
	SorcerySpeed bool
}

func (e *SorcerySpeed) Evaluate(ctx *evalstate.EvalState) bool {
	return judge.CanPlaySorcerySpeed(ctx.Game, ctx.PlayerID, nil) == e.SorcerySpeed
}

type Mode struct {
	Mode string
}

func (e *Mode) Evaluate(ctx *evalstate.EvalState) bool {
	return ctx.Mode == e.Mode
}

func compare(gotAny any, op Op, wantAny any) bool {
	switch got := gotAny.(type) {
	case int:
		want, ok := wantAny.(int)
		if !ok {
			return false
		}
		switch op {
		case OpEqual:
			return got == want
		case OpGreaterThan:
			return got > want
		case OpGreaterThanOrEqual:
			return got >= want
		case OpLessThan:
			return got < want
		case OpLessThanOrEqual:
			return got <= want
		case OpNotEqual:
			return got != want
		default:
			return false
		}
	}
	return false
}
