package target

import "deckronomicon/packages/game/mtg"

type Target struct {
	Type mtg.TargetType
	ID   string
}

type TargetSpec interface {
	TargetType() mtg.TargetType
	Name() string
}

type NoneTargetSpec struct {
}

func (t NoneTargetSpec) Name() string {
	return string(mtg.TargetTypeNone)
}

func (t NoneTargetSpec) TargetType() mtg.TargetType {
	return mtg.TargetTypeNone
}

type PlayerTargetSpec struct {
}

func (t PlayerTargetSpec) Name() string {
	return string(mtg.TargetTypePlayer)
}

func (t PlayerTargetSpec) TargetType() mtg.TargetType {
	return mtg.TargetTypePlayer
}

type PermanentTargetSpec struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func (t PermanentTargetSpec) Name() string {
	return string(mtg.TargetTypePermanent)
}

func (t PermanentTargetSpec) TargetType() mtg.TargetType {
	return mtg.TargetTypePermanent
}

type SpellTargetSpec struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func (t SpellTargetSpec) Name() string {
	return string(mtg.TargetTypeSpell)
}

func (t SpellTargetSpec) TargetType() mtg.TargetType {
	return mtg.TargetTypeSpell
}

type CardTargetSpec struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}

func (t CardTargetSpec) Name() string {
	return string(mtg.TargetTypeCard)
}

func (t CardTargetSpec) TargetType() mtg.TargetType {
	return mtg.TargetTypeCard
}
