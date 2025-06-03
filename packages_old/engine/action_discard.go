package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/game/player"
	"fmt"
)

type DiscardToHandSizeAction struct {
	Player *player.Player
}

func (a *DiscardToHandSizeAction) Description() string {
	return fmt.Sprintf(
		"player '%s' discards down to %d cards",
		a.Player.ID(),
		a.Player.Hand().Size(),
	)
}

func (a *DiscardToHandSizeAction) RequiresChoices() bool {
	return a.Player.Hand().Size() > a.Player.MaxHandSize
}

func (a *DiscardToHandSizeAction) GetChoices() ([]choose.Choice, error) {
	return nil, nil
}

func (a *DiscardToHandSizeAction) Complete() error {
	return nil
}
