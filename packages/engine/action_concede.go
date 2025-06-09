package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ConcedeAction struct {
	player state.Player
}

func NewConcedeAction(player state.Player) ConcedeAction {
	return ConcedeAction{
		player: player,
	}
}

func (a ConcedeAction) PlayerID() string {
	return a.player.ID()
}

func (a ConcedeAction) Name() string {
	return "Concede"
}

func (a ConcedeAction) Description() string {
	return "The active player concedes the game."
}

func (a ConcedeAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Conceding the game",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ConcedeAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.ConcedeEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
