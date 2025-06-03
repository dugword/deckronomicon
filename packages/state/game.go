package engine

import "deckronomicon/packages/game/gob"

type GameState struct {
	activePlayerIdx int
	battlefield     []*gob.Card
	phase           string
	players         []*Player
	stack           Stack
	turn            int
}

// Returns the players starting with the active player and going in turn order
func (s *GameState) PlayersInTurnOrder() []string {
	var n = s.activePlayerIdx
	var playersInTurnOrder []string
	for i := 0; i < len(s.players); i++ {
		playersInTurnOrder = append(playersInTurnOrder, s.players[n].id)
		n = (n + 1) % len(s.players)
	}
	return playersInTurnOrder
}

func (s *GameState) ActivePlayer() *Player {
	return s.players[s.activePlayerIdx]
}

func (s *GameState) GetPlayer(id string) (*Player, error) {
	for _, player := range s.players {
		if player.id == id {
			return player, nil
		}
	}
	return nil, ErrPlayerNotFound
}

func (s *GameState) NextPlayer() *Player {
	s.activePlayerIdx = (s.activePlayerIdx + 1) % len(s.players)
	return s.players[s.activePlayerIdx]
}

func NewGameState(config GameStateConfig) *GameState {
	state := GameState{
		players: config.Players,
	}
	return &state
}

type GameStateConfig struct {
	Players []*Player
}

func (g *GameState) IsGameOver() bool {
	return false // TODO: real logic
}

type GameStateSnapshot struct {
	Turn int
}
