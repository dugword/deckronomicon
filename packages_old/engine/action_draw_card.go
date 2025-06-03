package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/player"
	"fmt"
)

type DrawCardAction struct {
	Player *player.Player
}

func (a *DrawCardAction) Description() string {
	return a.Player.ID() + " draws a card."
}

func (a *DrawCardAction) RequiresChoice() bool {
	return false
}

func (a *DrawCardAction) GetChoices(state *GameState) ([]choose.ChoicePrompt, error) {
	return nil, nil
}

func (a *DrawCardAction) Resolve(state *GameState, _ []choose.ChoiceResponse) (string, error) {
	cardName, err := a.Player.DrawCard()
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%s drew a card: %s", a.Player.ID(), cardName), nil
}
