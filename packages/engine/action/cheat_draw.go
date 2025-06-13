package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

type DrawCheatAction struct {
	player state.Player
}

func NewDrawCheatAction(player state.Player) DrawCheatAction {
	return DrawCheatAction{
		player: player,
	}
}

func (a DrawCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a DrawCheatAction) Name() string {
	return "Draw a Card"
}

func (a DrawCheatAction) Description() string {
	return "Draw a card from your hand."
}

// TODO: Refactor so we don't need GetPrompt
func (a DrawCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{}, nil
}

func (a DrawCheatAction) Complete(
	game state.Game,
	choiceResults choose.ChoiceResults,
) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	return []event.GameEvent{
		event.CheatDrawEvent{
			PlayerID: a.player.ID(),
		},
		event.DrawCardEvent{
			PlayerID: a.player.ID(),
		},
	}, nil
}
