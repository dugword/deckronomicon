package state

import (
	"errors"
)

type Game struct {
	activePlayerIdx       int
	playerWithPriority    string
	playersPassedPriority map[string]bool
	battlefield           Battlefield
	phase                 string
	players               []Player
	stack                 Stack
	turn                  int
	winnerID              string
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

// Returns the players starting with the active player and going in turn order
func (s Game) PlayerIDsInTurnOrder() []string {
	var n = s.activePlayerIdx
	var playersInTurnOrder []string
	for i := 0; i < len(s.players); i++ {
		playersInTurnOrder = append(playersInTurnOrder, s.players[n].id)
		n = (n + 1) % len(s.players)
	}
	return playersInTurnOrder
}

func (s Game) ActivePlayerID() string {
	return s.players[s.activePlayerIdx].ID()
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
func (s Game) GetPlayer(id string) (Player, error) {
	for _, player := range s.players {
		if player.id == id {
			return player, nil
		}
	}
	return Player{}, errors.New("player not found")
}

func (s Game) NextPlayerID(currentPlayerID string) string {
	for i, player := range s.players {
		if player.id == currentPlayerID {
			nextIdx := (i + 1) % len(s.players)
			return s.players[nextIdx].ID()
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
