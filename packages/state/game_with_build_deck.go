package state

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"fmt"
)

// TODO Not sure if I like this here.
func (g *Game) WithBuildDeck(
	deckList *configs.DeckList,
	cardDefinitions map[string]*definition.Card,
	playerID string,
) (*Game, []*gob.Card, error) {
	// TODO: Make this nicer
	newGameCopy := *g
	newGame := &newGameCopy
	var deck []*gob.Card
	for _, entry := range deckList.Cards {
		for range entry.Count {
			cardDefinition, ok := cardDefinitions[entry.Name]
			if !ok {
				return nil, nil, fmt.Errorf(
					"card %s not found in card definitions", entry.Name,
				)
			}
			var id string
			id, newGame = newGame.WithGetNextID()
			cardDefinition.ID = id
			cardDefinition.Controller = playerID
			cardDefinition.Owner = playerID
			c := gob.NewCardFromDefinition(cardDefinition)
			deck = append(deck, c)
		}
	}
	return newGame, deck, nil
}
