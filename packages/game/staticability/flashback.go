package staticability

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
)

type Flashback struct {
	Cost cost.Cost
}

func (a *Flashback) Name() string {
	return "Flashback"
}

func (a *Flashback) StaticKeyword() mtg.StaticKeyword {
	return mtg.StaticKeywordFlashback
}
