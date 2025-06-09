package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type ConjureCardCheatAction struct {
	player   state.Player
	cardName string
}

func NewConjureCardCheatAction(player state.Player, cardName string) ConjureCardCheatAction {
	return ConjureCardCheatAction{
		player:   player,
		cardName: cardName,
	}
}

func (a ConjureCardCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a ConjureCardCheatAction) Name() string {
	return "Conjure Card"
}

func (a ConjureCardCheatAction) Description() string {
	return "Conjure a card into your hand."
}

func (a ConjureCardCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Conjure a card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a ConjureCardCheatAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	return []event.GameEvent{event.NoOpEvent{
		Message: "Conjured a card into your hand",
	}}, nil
}
