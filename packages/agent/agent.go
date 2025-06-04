package agent

import (
	"deckronomicon/packages/engine"
	"deckronomicon/packages/state"
	"fmt"
	"strings"
)

type Agent struct {
	id string
}

func NewAgent(id string) *Agent {
	agent := Agent{id: id}
	return &agent
}

func (a *Agent) GetNextAction() (engine.Action, error) {
	// press enter to continue
	fmt.Println("Press Enter to continue for agent:", a.id)
	fmt.Scanln() // wait for Enter Key
	return engine.PassAction{}, nil
}

func (a *Agent) ReportState(game state.Game) error {
	// For now, just print the game state
	fmt.Println("Reporting state for agent:", a.id)
	player, err := game.GetPlayer(a.id)
	if err != nil {
		return fmt.Errorf("report state: %w", err)
	}
	var cardNames []string
	for _, card := range player.Hand().GetAll() {
		cardNames = append(cardNames, card.Name())
	}
	fmt.Printf("Player %s: %s\n", player.ID(), strings.Join(cardNames, ", "))
	return nil
}
