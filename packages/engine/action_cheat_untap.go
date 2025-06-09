package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
)

type UntapCheatAction struct {
	player    state.Player
	permanent gob.Permanent
}

func NewUntapCheatAction(player state.Player, permanent gob.Permanent) UntapCheatAction {
	return UntapCheatAction{
		player:    player,
		permanent: permanent,
	}
}

func (a UntapCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a UntapCheatAction) Name() string {
	return "Untap target permanent"
}

func (a UntapCheatAction) Description() string {
	return "Untap target permanent."
}

func (a UntapCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Untap Permanent",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a UntapCheatAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Untap target permanent",
	}}, nil
}
