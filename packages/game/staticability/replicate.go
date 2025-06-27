package staticability

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/mtg"
)

type Replicate struct {
	Cost cost.Cost
}

func (a Replicate) Name() string {
	return "Replicate"
}

func (a Replicate) StaticKeyword() mtg.StaticKeyword {
	return mtg.StaticKeywordReplicate
}
