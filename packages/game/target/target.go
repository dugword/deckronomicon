package target

import "deckronomicon/packages/query"

// TODO Maybe target is its own package.
// Could be a compound type like cost

// TODO: This maybe should live in game/target

// This could be a thing with json.RawMessage like the effects.
// Then I could still have type safety with targets, and have different target types.
// and still have it be JSON serializable.
// TODO: Either these could be slice values, or I could have TargetValue be a slice.
type TargetValue struct {
	TargetType TargetType `json:"TargetType"`
	PlayerID   string     `json:"PlayerID,omitempty"`
	ObjectID   string     `json:"ObjectID,omitempty"`
}

type TargetSpec interface {
	TargetType() TargetType
	Name() string
}

type TargetType string

const (
	TargetTypeNone      TargetType = "None"
	TargetTypePlayer    TargetType = "Player"
	TargetTypeCreature  TargetType = "Creature"
	TargetTypeSpell     TargetType = "Spell"
	TargetTypePermanent TargetType = "Permanent"
)

type NoneTargetSpec struct {
}

func (t NoneTargetSpec) Name() string {
	return string(TargetTypeNone)
}

func (t NoneTargetSpec) TargetType() TargetType {
	return TargetTypeNone
}

type PlayerTargetSpec struct {
}

func (t PlayerTargetSpec) Name() string {
	return string(TargetTypePlayer)
}

func (t PlayerTargetSpec) TargetType() TargetType {
	return TargetTypePlayer
}

type PermanentTargetSpec struct{}

func (t PermanentTargetSpec) Name() string {
	return string(TargetTypePermanent)
}

func (t PermanentTargetSpec) TargetType() TargetType {
	return TargetTypePermanent
}

type SpellTargetSpec struct {
	Predicate query.Predicate
}

func (t SpellTargetSpec) Name() string {
	return string(TargetTypeSpell)
}

func (t SpellTargetSpec) TargetType() TargetType {
	return TargetTypeSpell
}
