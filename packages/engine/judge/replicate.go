package judge

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/state"
	"fmt"
)

func CanReplicateCard(
	game *state.Game,
	playerID string,
	cardToReplicate *gob.Card,
	ruling *Ruling,
	maybeApply func(game *state.Game, event event.GameEvent) (*state.Game, error),
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
	replicateAbility, ok := staticAbility.(*staticability.Replicate)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q has replicate ability, but it is not a Replicate ability", cardToReplicate.ID()))
		}
		can = false
		return can
	}
	if !CanPayCost(replicateAbility.Cost, cardToReplicate, game, playerID, ruling, maybeApply) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player %q cannot pay cost %s for card %q", playerID, replicateAbility.Cost, cardToReplicate.ID()))
		}
		can = false
	}
	return can
}
