package gob

import "deckronomicon/packages/game/mtg"

// TODO: I don't think I like this here

type SpliceModifiers struct {
	Subtype mtg.Subtype `json:"Subtype,omitempty"`
}
