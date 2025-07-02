package evaluator

import (
	"deckronomicon/packages/agent/auto/strategy/evalstate"
	"deckronomicon/packages/engine/judge"
)

type ManaPoolStat string

const (
	ManaPoolWhite     ManaPoolStat = "White"
	ManaPoolBlue      ManaPoolStat = "Blue"
	ManaPoolBlack     ManaPoolStat = "Black"
	ManaPoolRed       ManaPoolStat = "Red"
	ManaPoolGreen     ManaPoolStat = "Green"
	ManaPoolColorless ManaPoolStat = "Colorless"
	ManaPoolAny       ManaPoolStat = "Any"
)

type AvailableManaStat string

const (
	AvailableManaWhite     AvailableManaStat = "White"
	AvailableManaBlue      AvailableManaStat = "Blue"
	AvailableManaBlack     AvailableManaStat = "Black"
	AvailableManaRed       AvailableManaStat = "Red"
	AvailableManaGreen     AvailableManaStat = "Green"
	AvailableManaColorless AvailableManaStat = "Colorless"
	AvailableManaAny       AvailableManaStat = "Any"
)

type ManaPool struct {
	Stat  ManaPoolStat
	Op    Op
	Value any
}

func (e *ManaPool) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	var statValue any
	switch e.Stat {
	case ManaPoolWhite:
		statValue = player.ManaPool().White()
	case ManaPoolBlue:
		statValue = player.ManaPool().Blue()
	case ManaPoolBlack:
		statValue = player.ManaPool().Black()
	case ManaPoolRed:
		statValue = player.ManaPool().Red()
	case ManaPoolGreen:
		statValue = player.ManaPool().Green()
	case ManaPoolColorless:
		statValue = player.ManaPool().Colorless()
	case ManaPoolAny:
		statValue = player.ManaPool().Total()
	default:
		return false
	}
	return compare(statValue, e.Op, e.Value)
}

type AvailableMana struct {
	Stat  AvailableManaStat
	Op    Op
	Value any
}

func (e *AvailableMana) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	var statValue any
	availableMana := judge.GetAvailableMana(ctx.Game, player.ID())
	switch e.Stat {
	case AvailableManaWhite:
		statValue = availableMana.White()
	case AvailableManaBlue:
		statValue = availableMana.Blue()
	case AvailableManaBlack:
		statValue = availableMana.Black()
	case AvailableManaRed:
		statValue = availableMana.Red()
	case AvailableManaGreen:
		statValue = availableMana.Green()
	case AvailableManaColorless:
		statValue = availableMana.Colorless()
	case AvailableManaAny:
		statValue = availableMana.Total()
	default:
		return false
	}
	return compare(statValue, e.Op, e.Value)
}
