package dummy

import (
	"deckronomicon/game"
)

// TODO: This could be a rule-based agent where the rules are always to pass.

type DummyAgent struct {
	playerID string
}

func NewDummyAgent(playerID string) *DummyAgent {
	return &DummyAgent{playerID: playerID}
}

func (d *DummyAgent) ChooseMany(prompt string, source game.ChoiceSource, choices []game.Choice) ([]game.Choice, error) {
	return choices, nil // Always choose all available options
}

func (d *DummyAgent) ChooseOne(prompt string, source game.ChoiceSource, choices []game.Choice) (game.Choice, error) {
	return choices[0], nil // Always choose the first option
}

func (d *DummyAgent) Confirm(prompt string, source game.ChoiceSource) (bool, error) {
	return true, nil // Always confirm
}

func (d *DummyAgent) EnterNumber(prompt string, source game.ChoiceSource) (int, error) {
	return 1, nil // Always enter 1
}

func (d *DummyAgent) GetNextAction(state *game.GameState) *game.GameAction {
	return &game.GameAction{Type: game.ActionPass}
}

func (d *DummyAgent) PlayerID() string {
	return d.playerID
}

func (d *DummyAgent) ReportState(g *game.GameState) {}
