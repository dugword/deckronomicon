package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
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
	staticAbility, ok := cardToReplicate.StaticAbility(mtg.StaticKeywordReplicate)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q does not have replicate ability", cardToReplicate.ID()))
		}
		can = false
		return can
	}
	replicateAbility, ok := staticAbility.(staticability.Replicate)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q has replicate ability, but it is not a Replicate ability", cardToReplicate.ID()))
		}
		can = false
		return can
	}
	if !CanPayCost(replicateAbility.Cost, cardToReplicate, game, player, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player %q cannot pay cost %s for card %q", player.ID(), replicateAbility.Cost, cardToReplicate.ID()))
		}
		can = false
	}
	return can
}
