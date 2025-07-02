package engine

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type TurnBasedAction interface {
	Name() string
	Description() string
	GetPrompt(*state.Game) (choose.ChoicePrompt, error)
	Complete(*state.Game, choose.ChoiceResults) ([]event.GameEvent, error)
	PlayerID() string
}

var ErrInvalidTurnAction = errors.New("invalid turn action")

func (e *Engine) CompleteTurnAction(action TurnBasedAction) error {
	choicePrompt, err := action.GetPrompt(e.game)
	if err != nil {
		return fmt.Errorf(
			"failed to get choice prompt for action %q: %w",
			action.Name(),
			err,
		)
	}
	var choiceResults choose.ChoiceResults
	if choicePrompt.ChoiceOpts != nil {
		cs, err := e.agents[action.PlayerID()].Choose(choicePrompt)
		if err != nil {
			return fmt.Errorf(
				"failed to get choices for action %q: %w",
				action.Name(),
				err,
			)
		}
		choiceResults = cs
	}
	evnts, err := action.Complete(e.game, choiceResults)
	if err != nil {
		return fmt.Errorf(
			"failed to complete action %q: %w",
			action.Name(),
			errors.Join(err, ErrInvalidUserAction),
		)
	}
	for _, evnt := range evnts {
		if err := e.ApplyEvent(evnt); err != nil {
			return fmt.Errorf(
				"failed to apply event %q: %w",
				evnt.EventType(),
				err,
			)
		}
	}
	return nil
}
