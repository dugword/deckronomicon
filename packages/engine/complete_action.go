package engine

// TODO Document what level things happen at.

// Maybe Action.complete and Spell|Ability.resolve just takes the
// choices/targets and generates the events that need to be applied to the
// game.

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

type Action interface {
	Name() string
	Complete(state.Game) ([]event.GameEvent, error)
	Description() string
	PlayerID() string
}

var ErrInvalidUserAction = errors.New("invalid user action")

func (e *Engine) CompleteAction(action Action) error {
	evnts, err := action.Complete(e.game)
	if err != nil {
		return fmt.Errorf(
			"failed to complete action %q: %w",
			action.Name(),
			errors.Join(err, ErrInvalidUserAction),
		)
	}
	for _, evnt := range evnts {
		if err := e.Apply(evnt); err != nil {
			return fmt.Errorf(
				"failed to apply event %q: %w",
				evnt.EventType(),
				err,
			)
		}
	}
	return nil
}
