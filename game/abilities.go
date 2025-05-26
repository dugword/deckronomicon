package game

// TODO: Figure this out
// Abilities are Game Objects too....

// TODO: Is this useful?
// Ability is the general interface for all abilities.
/*
type Ability interface {
	Description() string
	Resolve(game *GameState, agent PlayerAgent) error
}
*/

const (
	SpeedInstant string = "Instant"
	SpeedSorcery string = "Sorcery"
)

const (
	AbilityTagDiscard     string = "Discard"
	AbilityTagDraw        string = "Draw"
	AbilityTagManaAbility string = "ManaAbility"
	AbilityTagScry        string = "Scry"
)
