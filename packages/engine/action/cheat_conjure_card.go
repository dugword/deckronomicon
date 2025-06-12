package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
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
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	if _, ok := env.Definitions[a.cardName]; !ok {
		return nil, fmt.Errorf("card %q not found in definitions", a.cardName)
	}
	return []event.GameEvent{event.CheatConjureCardEvent{
		PlayerID: a.player.ID(),
		CardName: a.cardName,
	}}, nil
}
