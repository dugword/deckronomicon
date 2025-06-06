package engine

// TODO Document what level things happen at.

// Maybe Action.complete and Spell|Ability.resolve just takes the
// choices/targets and generates the events that need to be applied to the
// game.

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

type Action interface {
	Name() string
	Description() string
	GetPrompt(state.Game) (choose.ChoicePrompt, error)
	Complete(state.Game, []choose.Choice) ([]event.GameEvent, error)
	PlayerID() string
}

type PlayCardAction struct {
	playerID string
	cardID   string
}

func NewPlayCardAction(playerID string, cardID string) PlayCardAction {
	return PlayCardAction{
		playerID: playerID,
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
		}}, nil
	}
	return []event.GameEvent{event.CastSpellEvent{
		PlayerID: a.playerID,
		CardID:   card.ID(),
	}}, nil
}

type PassPriorityAction struct {
	playerID string
}

func NewPassPriorityAction(playerID string) PassPriorityAction {
	return PassPriorityAction{
		playerID: playerID,
	}
}

func (a PassPriorityAction) PlayerID() string {
	return a.playerID
}

func (a PassPriorityAction) Name() string {
	return "Pass Priority"
}

func (a PassPriorityAction) Description() string {
	return "The active player passes priority to the next player."
}

func (a PassPriorityAction) GetPrompt(state state.Game) (choose.ChoicePrompt, error) {
	// No player choice needed, but we still return an empty prompt for consistency
	return choose.ChoicePrompt{
		Message:  "Passing priority",
		Choices:  nil,
		Optional: false,
	}, nil
}

func (a PassPriorityAction) Complete(
	game state.Game,
	choices []choose.Choice,
) ([]event.GameEvent, error) {
	playerID := game.PriorityPlayerID()
	return []event.GameEvent{event.PassPriorityEvent{
		// TODO: Need to think about how this is managed
		PlayerID: playerID,
	}}, nil
}
