package state

import (
	"deckronomicon/packages/game/mtg"
	"errors"
)

type Game struct {
	activePlayerIdx       int
	playerWithPriority    string
	playersPassedPriority map[string]bool
	battlefield           Battlefield
	phase                 mtg.Phase
	step                  mtg.Step
	players               []Player
	stack                 Stack
	winnerID              string
}

func (g Game) WithPhase(phase mtg.Phase) Game {
	g.phase = phase
	return g
}

func (g Game) WithStep(step mtg.Step) Game {
	g.step = step
	return g
}

func (g Game) Phase() mtg.Phase {
	return g.phase
}

func (g Game) Step() mtg.Step {
	return g.step
}

func (g Game) IsStackEmtpy() bool {
	return g.stack.Size() == 0
}

func (g Game) WithGameOver(winnerID string) Game {
	g.winnerID = winnerID
	return g
}

func (g Game) WithPlayers(players []Player) Game {
	g.players = players
	return g
}

func (g Game) WithClearedPriority() Game {
	g.playerWithPriority = ""
	return g
}

func (g Game) WithPlayerWithPriority(playerID string) Game {
	g.playerWithPriority = playerID
	return g
}

func (g Game) WithActivePlayer(playerID string) Game {
	var idx int
	for i, p := range g.players {
		if p.id == playerID {
			idx = i
			break
		}
	}
	g.activePlayerIdx = idx
	return g
}

func (g Game) WithResetPriorityPasses() Game {
	g.playersPassedPriority = map[string]bool{}
	return g
}

func (g Game) WithPlayerPassedPriority(playerID string) Game {
	newPlayersPassedPriority := map[string]bool{}
	for pID := range g.playersPassedPriority {
		newPlayersPassedPriority[pID] = g.playersPassedPriority[pID]
	}
	newPlayersPassedPriority[playerID] = true
	g.playersPassedPriority = newPlayersPassedPriority
	return g
}

func (g Game) WithUpdatedPlayer(player Player) Game {
	var newPlayers []Player
	for _, p := range g.players {
		if p.id == player.id {
			newPlayers = append(newPlayers, player)
			continue
		}
		newPlayers = append(newPlayers, p)
	}
	g.players = newPlayers
	return g
}

func (g Game) WithBattlefield(battlefield Battlefield) Game {
	g.battlefield = battlefield
	return g
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
func (g Game) GetPlayer(id string) (Player, error) {
	for _, player := range g.players {
		if player.id == id {
			return player, nil
		}
	}
	return Player{}, errors.New("player not found")
}

// TODO: THIS WILL BREAK WITH MORE THAN 2 PLAYERS
func (g Game) GetOpponent(id string) (Player, error) {
	if len(g.players) > 2 {
		panic("GetOpponent is not implemented for more than 2 players")
	}
	opponentID := g.NextPlayerID(id)
	for _, player := range g.players {
		if player.id == opponentID {
			return player, nil
		}
	}
	return Player{}, errors.New("opponent not found")
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
	if g.winnerID != "" {
		return true
	}
	return false
}

type GameStateSnapshot struct {
	Turn int
}
