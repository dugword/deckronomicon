package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ViewAction struct {
	player state.Player
	zone   string
	cardID string
}

func NewViewAction(player state.Player, zone string, cardID string) ViewAction {
	return ViewAction{
		player: player,
		zone:   zone,
		cardID: cardID,
	}
}

func (a ViewAction) PlayerID() string {
	return a.player.ID()
}

func (a ViewAction) Name() string {
	return "View card"
}

func (a ViewAction) Description() string {
	return "View a card in any zone."
}

func (a ViewAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "View card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ViewAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Viewed card in zone " + a.zone + ": " + a.cardID,
	}}, nil
}
