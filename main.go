package main

import (
	"deckronomicon/packages/agent"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
)

func main() {
	var seed int64 = 0
	fmt.Println("Starting the application...")
	fmt.Println("Initializing the engine...")
	player1ID := "Player1"
	player2ID := "Player2"
	deckLists := map[string][]gob.Card{
		player1ID: {
			gob.NewCard("id", "Plains"),
			gob.NewCard("id", "Island"),
			gob.NewCard("id", "Swamp"),
			gob.NewCard("id", "Mountain"),
			gob.NewCard("id", "Forest"),
		},
		player2ID: {
			gob.NewCard("id", "Plains"),
			gob.NewCard("id", "Island"),
			gob.NewCard("id", "Swamp"),
			gob.NewCard("id", "Mountain"),
			gob.NewCard("id", "Forest"),
		},
	}
	players := []state.Player{
		state.NewPlayer(player1ID, deckLists[player1ID]),
		state.NewPlayer(player2ID, deckLists[player2ID]),
	}
	agents := map[string]engine.PlayerAgent{
		player1ID: agent.NewAgent(player1ID),
		player2ID: agent.NewAgent(player2ID),
	}
	engineConfig := engine.EngineConfig{
		Agents:    agents,
		DeckLists: deckLists,
		Players:   players,
		Seed:      seed,
	}
	engine := engine.NewEngine(engineConfig)
	if err := engine.RunGame(); err != nil {
		fmt.Println("Error running the game:", err)
		return
	}
}
