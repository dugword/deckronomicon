package main

import (
	"bufio"
	"deckronomicon/packages/agent/dummy"
	"deckronomicon/packages/agent/interactive"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
	"os"
)

func generateDeckList(playerID string) []gob.Card {
	var deckList []gob.Card
	cardList := []string{
		"Plains",
		"Island",
		"Swamp",
		"Mountain",
		"Forest",
	}
	for i, card := range cardList {
		for j := range 12 {
			deckList = append(deckList, gob.NewCard(
				fmt.Sprintf("%d", (i+1)*(j+1)),
				card,
			))
		}
	}
	return deckList
}

func main() {
	stdin := os.Stdin
	var seed int64 = 0
	fmt.Println("Starting the application...")
	fmt.Println("Initializing the engine...")
	player1ID := "Player1"
	player2ID := "Player2"
	deckLists := map[string][]gob.Card{
		player1ID: generateDeckList(player1ID),
		player2ID: generateDeckList(player2ID),
	}
	players := []state.Player{
		state.NewPlayer(player1ID, deckLists[player1ID]),
		state.NewPlayer(player2ID, deckLists[player2ID]),
	}
	scanner := bufio.NewScanner(stdin)
	agents := map[string]engine.PlayerAgent{
		player1ID: interactive.NewAgent(
			scanner,
			player1ID,
			[]mtg.Step{mtg.StepPrecombatMain},
			true,
		),
		player2ID: dummy.NewAgent(player2ID, nil, false),
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
