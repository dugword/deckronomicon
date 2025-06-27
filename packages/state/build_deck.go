package state

import (
	"deckronomicon/packages/configs"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/gob"
	"fmt"
)

// TODO Not sure if I like this here.
func (g Game) WithBuildDeck(
	deckList configs.DeckList,
	cardDefinitions map[string]definition.Card,
	playerID string,
) (Game, []gob.Card, error) {
	var deck []gob.Card
	for _, entry := range deckList.Cards {
		for range entry.Count {
			cardDefinition, ok := cardDefinitions[entry.Name]
			if !ok {
				return Game{}, nil, fmt.Errorf(
					"card %s not found in card definitions", entry.Name,
				)
			}
			var id string
			id, g = g.GetNextID()
			cardDefinition.ID = id
			cardDefinition.Controller = playerID
			cardDefinition.Owner = playerID
			c := gob.NewCardFromDefinition(cardDefinition)
			deck = append(deck, c)
		}
	}
	return g, deck, nil
}
