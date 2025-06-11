package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

type ResetLandDropCheatAction struct {
	player state.Player
}

func NewResetLandDropCheatAction(player state.Player) ResetLandDropCheatAction {
	return ResetLandDropCheatAction{
		player: player,
	}
}

func (a ResetLandDropCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a ResetLandDropCheatAction) Name() string {
	return "Reset Land Drop"
}

func (a ResetLandDropCheatAction) Description() string {
	return "Reset the land drop for the turn."
}

func (a ResetLandDropCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Reset Land Drop",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ResetLandDropCheatAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{event.CheatResetLandDropEvent{
		PlayerID: a.PlayerID(),
	}}, nil
}
