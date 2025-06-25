package target

import "deckronomicon/packages/game/mtg"

// TODO Maybe target is its own package.
// Could be a compound type like cost

// TODO: This maybe should live in game/target

type TargetValue struct {
	TargetType TargetType
	TargetID   string
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
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func (t SpellTargetSpec) Name() string {
	return string(TargetTypeSpell)
}

func (t SpellTargetSpec) TargetType() TargetType {
	return TargetTypeSpell
}
