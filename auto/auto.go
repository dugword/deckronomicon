package auto

import (
	game "deckronomicon/game"
)

// AutoAgent implements PlayerAgent for automatic play
type AutoPlayerAgent struct{}

func NewAutoPlayerAgent() *AutoPlayerAgent {
	return &AutoPlayerAgent{}
}

func (a *AutoPlayerAgent) ReportState(state *game.GameState) {

}

func (a *AutoPlayerAgent) GetMoreInput(state *game.GameState, prompt game.OptionPrompt) game.Choice {

	return game.Choice{}
}

func (a *AutoPlayerAgent) GetNextAction(state *game.GameState) game.GameAction {
	if len(state.Hand) > 0 {
		// Just play the first card if any
		return game.GameAction{
			Type: game.ActionPlay,
		}
	}
	return game.GameAction{Type: game.ActionPass}
}
