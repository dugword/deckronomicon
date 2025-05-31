package effect

// TODO Make this require something so only GameState is passed in
type State any

// TODO Make this require something so only Player is passed in
type Player any

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	// Should this be named Name since it's not a unique ID?
	id string
	// TODO: Should this be named handler?
	Apply       func(State, Player) error
	description string
	tags        []Tag
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
func (e *Effect) Tags() []Tag {
	// Return the tags associated with the effect.
	return e.tags
}

// Tag represents a tag associated with an effect. These are used to
// define the effect and its properties. E.g. a Draw effect will have a tag
// Key of "Count" and a Value of the number of cards to draw.
type Tag struct {
	Key   string
	Value string
}
