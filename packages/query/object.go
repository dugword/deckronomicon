package query

import "deckronomicon/packages/game/mtg"

type Object interface {
	CardTypes() []mtg.CardType
	Colors() mtg.Colors
	Name() string
	// StaticAbilities() []mtg.StaticKeyword
	Subtypes() []mtg.Subtype
	Supertypes() []mtg.Supertype
}
