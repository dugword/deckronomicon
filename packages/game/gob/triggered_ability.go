package gob

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
)

type TriggeredAbility struct {
	ID         string
	SourceID   string
	SourceName string
	PlayerID   string
	Duration   mtg.Duration
	Effects    []effect.Effect
	Trigger    Trigger
	OneShot    bool
}

type Trigger struct {
	EventType string
	Filter    Filter
}

// todo this is common
type Filter struct {
	CardTypes  []mtg.CardType
	Colors     []mtg.Color
	Subtypes   []mtg.Subtype
	ManaValues []int
}
