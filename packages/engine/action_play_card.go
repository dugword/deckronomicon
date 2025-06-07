package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayCardAction struct {
	playerID string
	cardID   string
	zone     mtg.Zone
}

func NewPlayCardAction(playerID string, zone mtg.Zone, cardID string) PlayCardAction {
	return PlayCardAction{
		playerID: playerID,
		cardID:   cardID,
		zone:     zone,
	}

}

func (a PlayCardAction) PlayerID() string {
	return a.playerID
}

func (a PlayCardAction) Name() string {
	return "Play Card"
}
func (a PlayCardAction) Description() string {
	return "The active player plays a card from their hand."
}
func (a PlayCardAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	/*
		choices := state.GetPlayableCards(a.playerID)
		if len(choices) == 0 {
			return choose.ChoicePrompt{
				Message:  "No playable cards",
				Choices:  nil,
				Optional: false,
			}, nil
		}

		return choose.ChoicePrompt{
			Message:  "Choose a card to play",
			Choices:  choices,
			Optional: false,
		}, nil
	*/
	return choose.ChoicePrompt{}, nil
}
func (a PlayCardAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	player, err := game.GetPlayer(a.playerID)
	if err != nil {
		return nil, err // Player not found
	}
	card, ok := player.Hand().Get(a.cardID)
	if !ok {
		return nil, fmt.Errorf("card with ID %s not found in hand", a.cardID)
	}
	if card.Match(is.Land()) {
		return []event.GameEvent{event.PlayLandEvent{
			PlayerID: a.playerID,
			CardID:   card.ID(),
			Zone:     a.zone,
		}}, nil
	}
	return []event.GameEvent{event.CastSpellEvent{
		PlayerID: a.playerID,
		CardID:   card.ID(),
		Zone:     a.zone,
	}}, nil
}
