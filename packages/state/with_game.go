package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
)

func (g Game) WithCheatsEnabled(enabled bool) Game {
	g.cheatsEnabled = enabled
	return g
}

func (g Game) WithPhase(phase mtg.Phase) Game {
	g.phase = phase
	return g
}

func (g Game) WithStep(step mtg.Step) Game {
	g.step = step
	return g
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

func (g Game) WithPutCardOnBattlefield(card gob.Card, playerID string) (Game, error) {
	id, newGame := g.GetNextID()
	permanent, err := gob.NewPermanent(id, card, playerID)
	if err != nil {
		return newGame, err
	}
	newBattlefield := newGame.battlefield.Add(permanent)
	newerGame := newGame.WithBattlefield(newBattlefield)
	return newerGame, nil
}
