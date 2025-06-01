package query

import "deckronomicon/packages/game/mtg"

type Object interface {
	// StaticAbilities() []mtg.StaticKeyword
	CardTypes() []mtg.CardType
	Colors() mtg.Colors
	Description() string
	ID() string
	Match(Predicate) bool
	Name() string
	Subtypes() []mtg.Subtype
	Supertypes() []mtg.Supertype
}
