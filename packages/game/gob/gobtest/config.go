package gobtest

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
)

type SpellConfig struct {
	ID         string
	Name       string
	Controller string
	Owner      string
	ManaCost   cost.ManaCost
	Target     *target.TargetValue
}

type PermanentConfig struct {
	ID                 string
	Controller         string
	Owner              string
	Name               string
	ActivatedAbilities []definition.ActivatedAbilitySpec
	Card               CardConfig
	Tapped             bool
	CardTypes          []mtg.CardType
}

// This and definition.Card could be one in the same.
// We could have a definition for all states and gobs.
// And be ablet to load and serialize them.
type CardConfig struct {
	ID                 string
	Name               string
	CardTypes          []mtg.CardType
	Description        string
	ManaCost           string
	Loyalty            int
	SpellAbility       []definition.EffectSpec
	RulesText          string
	Power              int
	Subtypes           []mtg.Subtype
	StaticAbilities    []definition.StaticAbilitySpec
	ActivatedAbilities []definition.ActivatedAbilitySpec
}
