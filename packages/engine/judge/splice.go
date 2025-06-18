package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
)

func CanSpliceCard(
	player state.Player,
	cardToCast gob.Card,
	spliceCard gob.Card,
	ruling *Ruling,
) bool {
	can := true
	if _, ok := player.GetCardFromZone(spliceCard.ID(), mtg.ZoneHand); !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "splice card is not in hand")
		}
		can = false
	}
	rawModifers, ok := spliceCard.GetStaticAbilityModifiers(mtg.StaticKeywordSplice)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q does not have splice ability", spliceCard.ID()))
		}
		can = false
		return can
	}
	var spliceModifers gob.SpliceModifiers
	if err := json.Unmarshal(rawModifers, &spliceModifers); err != nil {
		panic(fmt.Sprintf("failed to unmarshal splice modifiers for card %q: %s", spliceCard.ID(), err))
	}
	if !cardToCast.Match(has.Subtype(spliceModifers.Subtype)) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf(
				"card %q does not have subtype %q",
				cardToCast.ID(),
				spliceModifers.Subtype,
			))
		}
		can = false
	}
	return can
}
