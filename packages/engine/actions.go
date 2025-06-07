package engine

// TODO Document what level things happen at.

// Maybe Action.complete and Spell|Ability.resolve just takes the
// choices/targets and generates the events that need to be applied to the
// game.

import (
	"deckronomicon/packages/choose"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

type Action interface {
	Name() string
	Description() string
	GetPrompt(state.Game) (choose.ChoicePrompt, error)
	Complete(state.Game, []choose.Choice) ([]event.GameEvent, error)
	PlayerID() string
}
