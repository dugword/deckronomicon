package dummy

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/state"
)

// ChooseMinimumAgent is a dummy agent that always chooses the minimum number of options available.
// When a choice is optional, it will not make a choice.
// When a choice is required, it will select the minimum number of options specified by the choice prompt.
// Options will be selected in the order they are presented.
// It can be used as the default player agent for opponents in games where the player is not interactive.
// This agent should always return valid choices and prevent the game from stalling.
type ChooseMinimumAgent struct {
	playerID string
}

func NewChooseMinimumAgent(playerID string) *ChooseMinimumAgent {
	agent := ChooseMinimumAgent{
		playerID: playerID,
	}
	return &agent
}

func (a *ChooseMinimumAgent) PlayerID() string {
	return a.playerID
}

func (a *ChooseMinimumAgent) GetNextAction(game *state.Game) (engine.Action, error) {
	return action.NewPassPriorityAction(), nil
}

func (a *ChooseMinimumAgent) Choose(game *state.Game, prompt choose.ChoicePrompt) (choose.ChoiceResults, error) {
	return choose.ChooseMinimal(prompt)
}

func (a *ChooseMinimumAgent) ReportState(game *state.Game) {}
