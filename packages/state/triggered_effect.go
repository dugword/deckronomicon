package state

import (
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
)

// TODO: Think through how this is named.
// I think it is the collection of all things that can be triggered.
// As new things enters play, if they have triggered abilities, they get added here.

type TriggeredEffectOG struct {
	ID               string
	TriggerCondition TriggerCondition
	EffectSpecs      []definition.EffectSpec
	Duration         mtg.Duration
}

type TriggeredEffect struct {
	ID       string
	PlayerID string
	Duration mtg.Duration
	Effect   []definition.EffectSpec
	//Source  query.Object
	Trigger Trigger
}

type TriggerCondition struct {
	Type   string
	Filter query.Predicate
}

type EffectToApply struct {
	Type   string             // "AddMana"
	Target target.TargetValue // The target of the effect
}

type Trigger struct {
	EventType string
	Filter    Filter
}

type Filter struct {
	CardTypes []mtg.CardType
	Subtypes  []mtg.Subtype
}
