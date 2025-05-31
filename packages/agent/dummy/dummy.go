package dummy

import (
	"deckronomicon/packages/choice"
	"deckronomicon/packages/game/action"
	"deckronomicon/packages/game/player"
)

// TODO: This could be a rule-based agent where the rules are always to pass.

type DummyAgent struct {
	player *player.Player
}

func NewDummyAgent() *DummyAgent {
	return &DummyAgent{}
}

func (d *DummyAgent) ChooseMany(prompt string, source choice.Source, choices []choice.Choice) ([]choice.Choice, error) {
	return choices, nil // Always choose all available options
}

func (d *DummyAgent) ChooseOne(prompt string, source choice.Source, choices []choice.Choice) (choice.Choice, error) {
	return choices[0], nil // Always choose the first option
}

func (d *DummyAgent) Confirm(prompt string, source choice.Source) (bool, error) {
	return true, nil // Always confirm
}

func (d *DummyAgent) EnterNumber(prompt string, source choice.Source) (int, error) {
	return 1, nil // Always enter 1
}

func (d *DummyAgent) GetNextAction(state player.GameState) (action.Action, error) {
	return action.Action{Type: action.ActionPass}, nil
}

func (d *DummyAgent) RegisterPlayer(player *player.Player) {
	d.player = player
}

func (d *DummyAgent) ReportState(g player.GameState) {}
