package state

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
)

func NewGameFromDefinition(definition *definition.Game) *Game {
	game := Game{
		nextID:                definition.NextID,
		cheatsEnabled:         definition.CheatsEnabled,
		playersPassedPriority: definition.PlayersPassedPriority,
		winnerID:              definition.WinnerID,
		battlefield:           NewBattlefield(),
		stack:                 NewStack(),
		runID:                 definition.RunID,
	}
	if definition.Battlefield != nil {
		game.battlefield = NewBattlefieldFromDefinition(definition.Battlefield)
	}
	if definition.Players != nil {
		game.players = NewPlayersFromDefinition(definition.Players)
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
	return &game
}

func NewPlayerFromDefinition(definition *definition.Player) *Player {
	manaAmount, err := mana.ParseManaString(definition.ManaPool)
	if err != nil {
		panic("failed to parse mana pool: " + err.Error())
	}
	player := Player{
		id:                 definition.ID,
		life:               definition.Life,
		landPlayedThisTurn: definition.LandPlayedThisTurn,
		maxHandSize:        definition.MaxHandSize,
		hand:               NewHand(),
		library:            NewLibrary([]*gob.Card{}),
		exile:              NewExile(),
		graveyard:          NewGraveyard(),
		revealed:           NewRevealed(),
		manaPool:           mana.Pool{}.WithAddAmount(manaAmount),
	}
	if definition.Hand != nil {
		player.hand = &Hand{cards: NewCardsFromDefinitions(definition.Hand.Cards)}
	}
	if definition.Library != nil {
		player.library = &Library{cards: NewCardsFromDefinitions(definition.Library.Cards)}
	}
	if definition.Graveyard != nil {
		player.graveyard = &Graveyard{cards: NewCardsFromDefinitions(definition.Graveyard.Cards)}
	}
	if definition.Exile != nil {
		player.exile = &Exile{cards: NewCardsFromDefinitions(definition.Exile.Cards)}
	}
	if definition.Revealed != nil {
		player.revealed = &Revealed{cards: NewCardsFromDefinitions(definition.Revealed.Cards)}
	}
	return &player
}

func NewPlayersFromDefinition(definitions []*definition.Player) []*Player {
	var players []*Player
	for _, definition := range definitions {
		players = append(players, NewPlayerFromDefinition(definition))
	}
	return players
}

func NewCardsFromDefinitions(definitions []*definition.Card) []*gob.Card {
	var cards []*gob.Card
	for _, definition := range definitions {
		card := gob.NewCardFromDefinition(definition)
		cards = append(cards, card)
	}
	return cards
}

func NewBattlefieldFromDefinition(definition *definition.Battlefield) *Battlefield {
	var permanents []*gob.Permanent
	for _, permanentDefinition := range definition.Permanents {
		permanent := gob.NewPermanentFromDefinition(
			permanentDefinition,
		)
		permanents = append(permanents, permanent)
	}
	return &Battlefield{
		permanents: permanents,
	}
}
