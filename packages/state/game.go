package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"fmt"
)

type Game struct {
	nextID                       int
	cheatsEnabled                bool
	activePlayerIdx              int
	playersPassedPriority        map[string]bool
	battlefield                  *Battlefield
	phase                        mtg.Phase
	step                         mtg.Step
	players                      []*Player
	stack                        *Stack
	winnerID                     string
	registeredTriggeredAbilities []gob.RegisteredTriggeredAbility
	continuousEffects            []gob.ContinuousEffect
	runID                        string
}

func (g *Game) RunID() string {
	return g.runID
}

func (g *Game) CheatsEnabled() bool {
	return g.cheatsEnabled
}

func (g *Game) Phase() mtg.Phase {
	return g.phase
}

func (g *Game) Players() []*Player {
	return g.players
}

func (g *Game) Step() mtg.Step {
	return g.step
}

func (g *Game) IsStackEmtpy() bool {
	return g.stack.Size() == 0
}

func (g *Game) Battlefield() *Battlefield {
	return g.battlefield
}

func (g *Game) Stack() *Stack {
	return g.stack
}

// Returns the players starting with the active player and going in turn order
func (g *Game) PlayerIDsInTurnOrder() []string {
	var n = g.activePlayerIdx
	var playersInTurnOrder []string
	for i := 0; i < len(g.players); i++ {
		playersInTurnOrder = append(playersInTurnOrder, g.players[n].id)
		n = (n + 1) % len(g.players)
	}
	return playersInTurnOrder
}

func (g *Game) ActivePlayerID() string {
	return g.players[g.activePlayerIdx].ID()
}

func (g *Game) DidPlayerPassPriority(playerID string) bool {
	return g.playersPassedPriority[playerID]
}

func (g *Game) DidAllPlayersPassPriority() bool {
	for _, player := range g.players {
		if !g.playersPassedPriority[player.id] {
			return false
		}
	}
	return true
}

func (g *Game) GetPlayer(id string) *Player {
	for _, player := range g.players {
		if player.id == id {
			return player
		}
	}
	panic(fmt.Sprintf("player %s not found in game", id))
}

// TODO: THIS WILL BREAK WITH MORE THAN 2 PLAYERS
func (g *Game) GetOpponent(id string) *Player {
	if len(g.players) > 2 {
		panic("GetOpponent is not implemented for more than 2 players")
	}
	opponentID := g.NextPlayerID(id)
	return g.GetPlayer(opponentID)
}

func (g *Game) NextPlayerID(currentPlayerID string) string {
	for i, player := range g.players {
		if player.id == currentPlayerID {
			nextIdx := (i + 1) % len(g.players)
			return g.players[nextIdx].ID()
		}
	}
	return currentPlayerID
}

func (g *Game) PriorityPlayerID() string {
	priorityPlayer := g.ActivePlayerID()
	if g.DidAllPlayersPassPriority() {
		return priorityPlayer
	}
	for g.DidPlayerPassPriority(priorityPlayer) {
		priorityPlayer = g.NextPlayerID(priorityPlayer)
	}
	return priorityPlayer
}

func (g *Game) PlayerPassedPriority(id string) bool {
	return g.playersPassedPriority[id]
}

func (g *Game) IsGameOver() bool {
	return g.winnerID != ""
}

func (g *Game) RegisteredTriggeredAbilities() []gob.RegisteredTriggeredAbility {
	return g.registeredTriggeredAbilities
}
