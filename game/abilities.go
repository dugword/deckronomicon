package game

// TODO: Implement mana abilities to not use the stack

// Ability is the general interface for all abilities.
type Ability interface {
	Description() string
	Resolve(game *GameState, choiceResolver ChoiceResolver) error
}

type AbilityTag struct {
	Key   string
	Value string
}

// ActivatedAbility represents abilities that require activation costs.
type ActivatedAbility struct {
	Cost          Cost
	Description   string
	Effect        func(*GameState, ChoiceResolver)
	IsManaAbility bool // maybe make this a generic set of tags
	Tags          []AbilityTag
}

// StaticAbility represents continuous effects.
type StaticAbility struct {
	Effect      func(*GameState, ChoiceResolver)
	Description string
}

// SpellAbility represents abilities on instant or sorcery spells.
type SpellAbility struct {
	Effect      func(*GameState, ChoiceResolver)
	Description string
}

// TriggeredAbility represents abilities that trigger on specific events.
type TriggeredAbility struct {
	Description      string
	Effect           func(*GameState, ChoiceResolver)
	TriggerCondition func(event Event) bool
}
