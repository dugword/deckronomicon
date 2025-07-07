package gob

import (
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/mtg"
)

type RegisteredTriggeredAbility struct {
	ID       string
	Name     string
	SourceID string
	PlayerID string
	Duration mtg.Duration
	Effects  []effect.Effect
	Trigger  Trigger
	OneShot  bool
}
