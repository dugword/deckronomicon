package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

func CanReplicateCard(
	game state.Game,
	player state.Player,
	cardToReplicate gob.Card,
	ruling *Ruling,
) bool {
	can := true
	cost, ok := cardToReplicate.GetStaticAbilityCost(mtg.StaticKeywordReplicate)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q does not have replicate ability", cardToReplicate.ID()))
		}
		can = false
		return can
	}
	if !CanPayCost(cost, cardToReplicate, game, player, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player %q cannot pay cost %s for card %q", player.ID(), cost, cardToReplicate.ID()))
		}
		can = false
	}
	return can
}
