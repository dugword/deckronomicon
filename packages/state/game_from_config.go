package state

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state/statetest"
)

// This is probably the same as what we are doing with load from definition for cards,
// we sould rename this and come up with a consistent pattern.
func LoadGameFromConfig(config statetest.GameConfig) Game {
	return Game{
		nextID:                config.NextID,
		cheatsEnabled:         config.CheatsEnabled,
		activePlayerIdx:       config.ActivePlayerIdx,
		playersPassedPriority: config.PlayersPassedPriority,
		battlefield:           LoadBattlefieldFromConfig(config.Battlefield),
		phase:                 config.Phase,
		step:                  config.Step,
		players:               LoadPlayersFromConfig(config.Players),
		winnerID:              config.WinnerID,
	}
}

func LoadPlayerFromConfig(config statetest.PlayerConfig) Player {
	manaAmount, err := mana.ParseManaString(config.ManaPool)
	if err != nil {
		panic("failed to parse mana pool: " + err.Error())
	}
	return Player{
		id:        config.ID,
		life:      config.Life,
		hand:      LoadHandFromConfig(config.ID, config.Hand),
		library:   LoadLibraryFromConfig(config.ID, config.Library),
		graveyard: LoadGraveyardFromConfig(config.ID, config.Graveyard),
		exile:     LoadExileFromConfig(config.ID, config.Exile),
		manaPool:  mana.NewManaPool().WithAddedAmount(manaAmount),
	}
}

func LoadPlayersFromConfig(configs []statetest.PlayerConfig) []Player {
	var players []Player
	for _, config := range configs {
		players = append(players, LoadPlayerFromConfig(config))
	}
	return players
}
func LoadHandFromConfig(playerID string, config statetest.HandConfig) Hand {
	var cards []gob.Card
	for _, cardConfig := range config.Cards {
		card := gob.LoadCardFromConfig(playerID, cardConfig)
		cards = append(cards, card)
	}
	return Hand{
		cards: cards,
	}
}

func LoadExileFromConfig(playerID string, config statetest.ExileConfig) Exile {
	var cards []gob.Card
	for _, cardConfig := range config.Cards {
		card := gob.LoadCardFromConfig(playerID, cardConfig)
		cards = append(cards, card)
	}
	return Exile{
		cards: cards,
	}
}

func LoadLibraryFromConfig(playerID string, config statetest.LibraryConfig) Library {
	var cards []gob.Card
	for _, cardConfig := range config.Cards {
		card := gob.LoadCardFromConfig(playerID, cardConfig)
		cards = append(cards, card)
	}
	return Library{
		cards: cards,
	}
}

func LoadGraveyardFromConfig(playerID string, config statetest.GraveyardConfig) Graveyard {
	var cards []gob.Card
	for _, cardConfig := range config.Cards {
		card := gob.LoadCardFromConfig(playerID, cardConfig)
		cards = append(cards, card)
	}
	return Graveyard{
		cards: cards,
	}
}

func LoadBattlefieldFromConfig(config statetest.BattlefieldConfig) Battlefield {
	var permanents []gob.Permanent
	for _, permanentConfigs := range config.Permanents {
		permanent := gob.LoadPermanentFromConfig(permanentConfigs)
		permanents = append(permanents, permanent)
	}
	return Battlefield{
		permanents: permanents,
	}
}
