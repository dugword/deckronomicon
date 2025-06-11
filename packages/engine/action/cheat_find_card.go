package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

type FindCardCheatAction struct {
	player   state.Player
	cardName string
}

func NewFindCardCheatAction(player state.Player, cardName string) FindCardCheatAction {
	return FindCardCheatAction{
		player:   player,
		cardName: cardName,
	}
}

func (a FindCardCheatAction) Name() string {
	return "Find Card"
}

func (a FindCardCheatAction) PlayerID() string {
	return a.player.ID()
}

func (a FindCardCheatAction) Description() string {
	return "Find a card into your hand."
}

func (a FindCardCheatAction) GetPrompt(game state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Find a card",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a FindCardCheatAction) Complete(
	game state.Game,
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	if !game.CheatsEnabled() {
		return nil, fmt.Errorf("no cheating you cheater")
	}
	if a.cardName == "" {
		return nil, fmt.Errorf("find card action requires a card name")
	}
	card, ok := a.player.Library().Find(has.Name(a.cardName))
	if !ok {
		return nil, fmt.Errorf("card %q not found in library", a.cardName)
	}

	return []event.GameEvent{
		event.CheatFindCardEvent{
			PlayerID: a.player.ID(),
			CardName: a.cardName,
		},
		event.MoveCardEvent{
			CardID:   card.ID(),
			FromZone: mtg.ZoneLibrary,
			ToZone:   mtg.ZoneHand,
			PlayerID: a.player.ID(),
		},
	}, nil
}
