package evaluator

import "deckronomicon/packages/agent/auto/strategy/evalstate"

type PlayerStat string

const (
	PlayerStatLife               PlayerStat = "Life"
	PlayerStatHandSize           PlayerStat = "HandSize"
	PlayerStatLibrarySize        PlayerStat = "LibrarySize"
	PlayerStatGraveyardSize      PlayerStat = "GraveyardSize"
	PlayerStatName               PlayerStat = "Name"
	PlayerStatLandPlayedThisTurn PlayerStat = "LandPlayedThisTurn"
	PlayerStatID                 PlayerStat = "ID"
	PlayerStatTurn               PlayerStat = "Turn"
)

type Player struct {
	Stat  PlayerStat
	Op    Op
	Value any
}

func (e *Player) Evaluate(ctx *evalstate.EvalState) bool {
	player := ctx.Game.GetPlayer(ctx.PlayerID)
	var statValue any
	switch e.Stat {
	case PlayerStatLife:
		statValue = player.Life()
	case PlayerStatHandSize:
		statValue = player.Hand().Size()
	case PlayerStatLibrarySize:
		statValue = player.Library().Size()
	case PlayerStatGraveyardSize:
		statValue = player.Graveyard().Size()
	case PlayerStatName:
		statValue = player.Name()
	case PlayerStatTurn:
		statValue = player.Turn()
	}
	return compare(statValue, e.Op, e.Value)
}
