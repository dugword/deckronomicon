package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/staticability"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"fmt"
)

func CanSpliceCard(
	game *state.Game,
	playerID string,
	cardToCast *gob.Card,
	spliceCard *gob.Card,
	ruling *Ruling,
) bool {
	can := true
	player := game.GetPlayer(playerID)
	if _, ok := player.GetCardFromZone(spliceCard.ID(), mtg.ZoneHand); !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "splice card is not in hand")
		}
		can = false
	}
	staticAbility, ok := spliceCard.StaticAbility(mtg.StaticKeywordSplice)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q does not have splice ability", spliceCard.ID()))
		}
		can = false
		return can
	}
	spliceAbility, ok := staticAbility.(*staticability.Splice)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q has splice ability, but it is not a Splice ability", spliceCard.ID()))
		}
		can = false
		return can
	}
	if !cardToCast.Match(has.Subtype(spliceAbility.Subtype)) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf(
				"card %q does not have subtype %q",
				cardToCast.ID(),
				spliceAbility.Subtype,
			))
		}
		can = false
	}
	return can
}
