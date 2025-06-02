package effect

import "deckronomicon/packages/game/core"

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	Apply func(core.State, core.Player) error
	// Should this be named Name since it's not a unique ID?
	id string
	// TODO: Should this be named handler?
	description string
	tags        []core.Tag
}

func (e *Effect) ID() string {
	// Return the ID as the name for now, but this could be improved
	// to return a more descriptive name based on the effect.
	return e.id
}

func (e *Effect) Description() string {
	// Return the description of the effect.
	return e.description
}

// Tags returns the tags associated with the effect.
func (e *Effect) Tags() []core.Tag {
	// Return the tags associated with the effect.
	return e.tags
}
