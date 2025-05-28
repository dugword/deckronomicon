package auto

import (
	"bufio"
	"deckronomicon/game"
	"fmt"
	"os"
	"strings"
)

// TODO: Split out the display from interactive agent and use that in both
// places.
func (a *RuleBasedAgent) debugPrintStuff(state *game.GameState) {
	player, err := state.GetPlayer(a.playerID)
	if err != nil {
		fmt.Println("ERROR: could not find player with ID =>", a.playerID)
		os.Exit(1)
	}
	fmt.Println("##########")
	fmt.Println("Mode:", player.Mode)
	fmt.Println("Turn:", player.Turn)
	fmt.Println("Step:", state.CurrentStep)
	fmt.Println("Life:", player.Life)
	fmt.Printf("Graveyard: %d cards\n", player.Graveyard.Size())
	fmt.Printf("Library: %d cards\n", player.Library.Size())
	fmt.Printf("Mana: %s\n", player.ManaPool.Describe())
	var permanentNames []string
	for _, permanent := range player.Battlefield.GetAll() {
		permanentNames = append(permanentNames, permanent.Name())
	}
	fmt.Println("Battlefield: ", strings.Join(permanentNames, ", "))
	var cardNames []string
	for _, card := range player.Hand.GetAll() {
		cardNames = append(cardNames, card.Name())
	}
	fmt.Println("Hand: ", strings.Join(cardNames, ", "))
	fmt.Println("Message:", state.Message)
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
