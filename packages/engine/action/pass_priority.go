package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type PassPriorityAction struct {
	player state.Player
}

func NewPassPriorityAction(player state.Player) PassPriorityAction {
	return PassPriorityAction{
		player: player,
	}
}

func (a PassPriorityAction) PlayerID() string {
	return a.player.ID()
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Description() string {
	return "The active player passes priority to the next player."
}

func (a PassPriorityAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a PassPriorityAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.PassPriorityEvent{
		PlayerID: a.player.ID(),
	}}, nil
}
