package gob

const TagManaAbility = "ManaAbility"

// Effect represents an effect that can be applied to a game state.
type Effect struct {
	name      string
	modifiers []Tag
	tags      []Tag
	// TODO: is optional just a modifer?
	optional bool
	// apply func(ore.State, core.Player) error
	// Should this be named Name since it's not a unique ID?

	// TODO: Should this be named handler?
	//description string
	//tags        []Tag

}

func (e Effect) Modifiers() []Tag {
	// Return the modifiers associated with the effect.
	return e.modifiers
}

// func NewEffect(id string, description string, tags []Tag) Effect {
func NewEffect(name string) Effect {
	effect := Effect{
		name: name,
	}
	return effect
}

func (e Effect) Name() string {
	// Return the ID as the name for now, but this could be improved
	// to return a more descriptive name based on the effect.
	return e.name
}

/*
func (e Effect) Apply(state core.State, player core.Player) error {
	return e.apply(state, player)
}
*/

/*
func (e Effect) Description() string {
	// Return the description of the effect.
	return e.description
}
*/
// Tags returns the tags associated with the effect.
func (e Effect) Tags() []Tag {
	// Return the tags associated with the effect.
	return e.tags
}
