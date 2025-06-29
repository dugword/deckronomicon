package state

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
)

func NewGameFromDefinition(definition definition.Game) Game {
	game := Game{
		nextID:                definition.NextID,
		cheatsEnabled:         definition.CheatsEnabled,
		playersPassedPriority: definition.PlayersPassedPriority,
		battlefield:           NewBattlefieldFromDefinition(definition.Battlefield),
		players:               NewPlayersFromDefinition(definition.Players),
		winnerID:              definition.WinnerID,
	}
	if definition.Phase != "" {
		phase, ok := mtg.StringToPhase(definition.Phase)
		if ok {
			panic("failed to parse phase: " + definition.Phase)
		}
		game.phase = phase
	}
	if definition.Step != "" {
		step, ok := mtg.StringToStep(definition.Step)
		if !ok {
			panic("failed to parse step: " + definition.Step)
		}
		game.step = step
	}
	for i, player := range game.players {
		if player.id == definition.ActivePlayerID {
			game.activePlayerIdx = i
			break
		}
	}
	return game
}

func NewPlayerFromDefinition(definition definition.Player) Player {
	manaAmount, err := mana.ParseManaString(definition.ManaPool)
	if err != nil {
		panic("failed to parse mana pool: " + err.Error())
	}
	return Player{
		id:                 definition.ID,
		life:               definition.Life,
		landPlayedThisTurn: definition.LandPlayedThisTurn,
		hand:               Hand{cards: NewCardsFromDefinitions(definition.Hand.Cards)},
		library:            Library{cards: NewCardsFromDefinitions(definition.Library.Cards)},
		graveyard:          Graveyard{cards: NewCardsFromDefinitions(definition.Graveyard.Cards)},
		exile:              Exile{cards: NewCardsFromDefinitions(definition.Exile.Cards)},
		manaPool:           mana.Pool{}.WithAddAmount(manaAmount),
	}
}

func NewPlayersFromDefinition(definitions []definition.Player) []Player {
	var players []Player
	for _, definition := range definitions {
		players = append(players, NewPlayerFromDefinition(definition))
	}
	return players
}

func NewCardsFromDefinitions(definitions []definition.Card) []gob.Card {
	var cards []gob.Card
	for _, definition := range definitions {
		card := gob.NewCardFromDefinition(definition)
		cards = append(cards, card)
	}
	return cards
}

func NewBattlefieldFromDefinition(definition definition.Battlefield) Battlefield {
	var permanents []gob.Permanent
	for _, permanentDefinition := range definition.Permanents {
		permanent := gob.NewPermanentFromDefinition(
			permanentDefinition,
		)
		permanents = append(permanents, permanent)
	}
	return Battlefield{
		permanents: permanents,
	}
}
