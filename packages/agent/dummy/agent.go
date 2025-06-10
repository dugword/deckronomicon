package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
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

func (a *Agent) PlayerID() string {
	return a.id
}

func (a *Agent) GetNextAction(game state.Game) (engine.Action, error) {
	player, ok := game.GetPlayer(a.id)
	if !ok {
		return nil, fmt.Errorf("player %q not found", a.id)
	}
	return engine.NewPassPriorityAction(player), nil
}

func (a *Agent) Choose(prompt choose.ChoicePrompt) ([]choose.Choice, error) {
	var selected []choose.Choice
	for i := range prompt.MinChoices {
		selected = append(selected, prompt.Choices[i])
	}
	return selected, nil
}

func (a *Agent) ReportState(game state.Game) error {
	return nil
}
