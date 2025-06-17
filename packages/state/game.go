package state

// TODO: Document what kind of logic lives where.

// I think the state should really just control the mechanics of movings
// things around, not checking the game rules. Games rule enforcement should
// happen in the engine

// stuff like can cast might even need to be moved to the engine package

import (
	"deckronomicon/packages/game/mtg"
	"strconv"
)

type Game struct {
	id                    int
	cheatsEnabled         bool
	activePlayerIdx       int
	playerWithPriority    string
	playersPassedPriority map[string]bool
	battlefield           Battlefield
	phase                 mtg.Phase
	step                  mtg.Step
	players               []Player
	stack                 Stack
	winnerID              string
	// TODO: Rename triggered abilities or something....
	triggeredEffects  []TriggeredEffect
	continuousEffects []ContinuousEffect
}

func (g Game) CheatsEnabled() bool {
	return g.cheatsEnabled
}

func (g Game) Phase() mtg.Phase {
	return g.phase
}

func (g Game) Players() []Player {
	return g.players
}

func (g Game) Step() mtg.Step {
	return g.step
}

func (g Game) IsStackEmtpy() bool {
	return g.stack.Size() == 0
}

func (g Game) Battlefield() Battlefield {
	return g.battlefield
}

func (g Game) Stack() Stack {
	return g.stack
}

// Returns the players starting with the active player and going in turn order
func (g Game) PlayerIDsInTurnOrder() []string {
	var n = g.activePlayerIdx
	var playersInTurnOrder []string
	for i := 0; i < len(g.players); i++ {
		playersInTurnOrder = append(playersInTurnOrder, g.players[n].id)
		n = (n + 1) % len(g.players)
	}
	return playersInTurnOrder
}

func (g Game) ActivePlayerID() string {
	return g.players[g.activePlayerIdx].ID()
}

func (g Game) AllPlayersPassedPriority() bool {
	for _, player := range g.players {
		if !g.playersPassedPriority[player.id] {
			return false
		}
	}
	return true
}

// TODO: Think about removing this and using GetPlayerID instead
func (g Game) GetPlayer(id string) (Player, bool) {
	for _, player := range g.players {
		if player.id == id {
			return player, true
		}
	}
	return Player{}, false
}

// TODO: THIS WILL BREAK WITH MORE THAN 2 PLAYERS
func (g Game) GetOpponent(id string) (Player, bool) {
	if len(g.players) > 2 {
		panic("GetOpponent is not implemented for more than 2 players")
	}
	opponentID := g.NextPlayerID(id)
	for _, player := range g.players {
		if player.id == opponentID {
			return player, true
		}
	}
	return Player{}, false
}

func (g Game) NextPlayerID(currentPlayerID string) string {
	for i, player := range g.players {
		if player.id == currentPlayerID {
			nextIdx := (i + 1) % len(g.players)
			return g.players[nextIdx].ID()
		}
	}
	return currentPlayerID
}

func (g Game) PriorityPlayerID() string {
	return g.playerWithPriority
}

func (g Game) PlayerPassedPriority(id string) bool {
	return g.playersPassedPriority[id]
}

func NewGame(config GameConfig) Game {
	state := Game{
		players: config.Players,
	}
	return state
}

type GameConfig struct {
	Players []Player
}

func (g Game) IsGameOver() bool {
	return g.winnerID != ""
}

type GameStateSnapshot struct {
	Turn int
}

func (g Game) GetNextID() (id string, game Game) {
	g.id++
	return strconv.Itoa(g.id), g
}

func (g Game) TriggeredEffects() []TriggeredEffect {
	return g.triggeredEffects
}
