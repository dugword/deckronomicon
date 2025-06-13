package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"fmt"
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
	return choose.ChoicePrompt{}, nil
}

func (a UntapCheatAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatUntapEvent{
			PlayerID: a.player.ID(),
		},
		event.UntapPermanentEvent{
			PlayerID:    a.player.ID(),
			PermanentID: a.permanent.ID(),
		},
	}, nil
}
