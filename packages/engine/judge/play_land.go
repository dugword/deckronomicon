package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: Should judge check if the card is in the player's hand?
func CanPlayLand(
	game state.Game,
	player state.Player,
	zone mtg.Zone,
	card gob.Card,
	ruling *Ruling,
) bool {
	can := true
	if !card.Match(is.Land()) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "card is not a land")
		}
		can = false
	}
	if zone != mtg.ZoneHand {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "zone is not hand")
		}
		can = false
	}
	if !player.ZoneContains(zone, has.ID(card.ID())) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player does not have card in zone %q", zone))
		}
		can = false
	}
	if player.LandPlayedThisTurn() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "land played this turn already")
		}
		can = false
	}
	if !CanPlaySorcerySpeed(game, player.ID(), ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot play land at instant speed")
		}
		can = false
	}
	return can
}
