package action

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/mana"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type PlayCardAction struct {
	player          state.Player
	cardInZone      gob.CardInZone
	ManaPayment     mana.Pool
	AdditionalCosts []string
	Targets         []string
}

func NewPlayCardAction(player state.Player, cardInZone gob.CardInZone) PlayCardAction {
	return PlayCardAction{
		player:     player,
		cardInZone: cardInZone,
	}
}

func (a PlayCardAction) PlayerID() string {
	return a.player.ID()
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
	env *ResolutionEnvironment,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	zone, ok := a.player.GetZone(a.cardInZone.Zone())
	if !ok {
		return nil, fmt.Errorf("zone %q not valid", a.cardInZone.Zone())
	}
	if !zone.Contains(has.ID(a.cardInZone.ID())) {
		return nil, fmt.Errorf(
			"card %q not found in zone %q",
			a.cardInZone.ID(),
			a.cardInZone.Zone(),
		)
	}
	if a.cardInZone.Card().Match(is.Land()) {
		return []event.GameEvent{
			event.PlayLandEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				Zone:     a.cardInZone.Zone(),
			},
			event.PutCardOnBattlefieldEvent{
				PlayerID: a.player.ID(),
				CardID:   a.cardInZone.ID(),
				FromZone: a.cardInZone.Zone(),
			},
		}, nil
	}
	return []event.GameEvent{event.CastSpellEvent{
		PlayerID: a.player.ID(),
		CardID:   a.cardInZone.ID(),
		Zone:     a.cardInZone.Zone(),
	}}, nil
}
