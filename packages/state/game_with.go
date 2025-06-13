package state

import (
	"deckronomicon/packages/game/definition"
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
	playersPassedPriority := map[string]bool{}
	for pID := range g.playersPassedPriority {
		playersPassedPriority[pID] = g.playersPassedPriority[pID]
	}
	playersPassedPriority[playerID] = true
	g.playersPassedPriority = playersPassedPriority
	return g
}

func (g Game) WithUpdatedPlayer(player Player) Game {
	var players []Player
	for _, p := range g.players {
		if p.id == player.id {
			players = append(players, player)
			continue
		}
		players = append(players, p)
	}
	g.players = players
	return g
}

func (g Game) WithBattlefield(battlefield Battlefield) Game {
	g.battlefield = battlefield
	return g
}

// TODO: Should the permanent be created here or in the apply reducer...
func (g Game) WithPutPermanentOnBattlefield(card gob.Card, playerID string) (Game, error) {
	id, game := g.GetNextID()
	permanent, err := gob.NewPermanent(id, card, playerID)
	if err != nil {
		return game, err
	}
	battlefield := game.battlefield.Add(permanent)
	game = game.WithBattlefield(battlefield)
	return game, nil
}

func (g Game) WithStack(stack Stack) Game {
	g.stack = stack
	return g
}

func (g Game) WithPutSpellOnStack(card gob.Card, playerID string, flashback bool) (Game, error) {
	id, game := g.GetNextID()
	spell, err := gob.NewSpell(id, card, playerID, flashback)
	if err != nil {
		return game, err
	}
	stack := game.stack.AddTop(spell)
	game = game.WithStack(stack)
	return game, nil
}

func (g Game) WithPutAbilityOnStack(
	playerID,
	sourceID,
	abilityID,
	abilityName string,
	effectSpecs []definition.EffectSpec,
) (Game, error) {
	id, game := g.GetNextID()
	abilityOnStack := gob.NewAbilityOnStack(
		id,
		playerID,
		sourceID,
		abilityID,
		abilityName,
		effectSpecs,
	)
	stack := game.stack.AddTop(abilityOnStack)
	game = game.WithStack(stack)
	return game, nil
}
