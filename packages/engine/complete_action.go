package engine

// TODO Document what level things happen at.

// Maybe Action.complete and Spell|Ability.resolve just takes the
// choices/targets and generates the events that need to be applied to the
// game.

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/action"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

type Action interface {
	Name() string
	//Description() string
	// TODO: Don't prompt, pass in the choices directly. Need to figure out mulitgans and end of turn of
	// discard. I guess also declare attackers, declare blockers, and assign combat damage.
	// Maybe I need to split out Turn Based Actions from Player Actions.
	GetPrompt(state.Game) (choose.ChoicePrompt, error)
	Complete(state.Game, *action.ResolutionEnvironment, []choose.Choice) ([]event.GameEvent, error)
	PlayerID() string
}

func (e *Engine) CompleteAction(action Action) error {
	choicePrompt, err := action.GetPrompt(e.game)
	if err != nil {
		return fmt.Errorf(
			"failed to get choice prompt for action %q: %w",
			action.Name(),
			err,
		)
	}
	choices := []choose.Choice{}
	if len(choicePrompt.Choices) != 0 {
		cs, err := e.agents[action.PlayerID()].Choose(choicePrompt)
		if err != nil {
			return fmt.Errorf(
				"failed to get choices for action %q: %w",
				action.Name(),
				err,
			)
		}
		choices = cs
	}
	evnts, err := action.Complete(e.game, e.resolutionEnvironment, choices)
	if err != nil {
		return fmt.Errorf(
			"failed to complete action %q: %w",
			action.Name(),
			err,
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
