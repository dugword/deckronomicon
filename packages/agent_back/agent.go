package agent

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
	"strings"
)

type Agent struct {
	id      string
	verbose bool
	stops   []mtg.Step
}

func NewAgent(id string, stops []mtg.Step, verbose bool) *Agent {
	agent := Agent{
		id:      id,
		stops:   stops,
		verbose: verbose,
	}
	return &agent
}

func (a *Agent) GetNextAction(game state.Game) (engine.Action, error) {
	pass := true
	for _, step := range a.stops {
		if game.Step() == step {
			pass = false
			break
		}
	}
	if pass {
		return engine.NewPassAction(a.id), nil
	}

	// press enter to continue
	fmt.Println("Press Enter to continue for agent:", a.id)
	fmt.Scanln() // wait for Enter Key
	return engine.NewPassAction(a.id), nil
}

func (a *Agent) Choose(prompt choose.ChoicePrompt) ([]choose.Choice, error) {
	fmt.Println(prompt.Message)
	for i, choice := range prompt.Choices {
		fmt.Printf("%d: %s\n", i+1, choice.Name)
	}
	if prompt.Choices != nil {
		fmt.Println("Press Enter to continue for agent:", a.id)
		fmt.Scanln() // wait for Enter Key
	}
	return []choose.Choice{}, nil
}

func (a *Agent) ReportState(game state.Game) error {
	if !a.verbose {
		return nil
	}
	if game.Step() != mtg.StepPrecombatMain {
		return nil
	}
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
