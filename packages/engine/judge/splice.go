package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
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
	modifiers, ok := spliceCard.GetStaticAbilityModifiers(mtg.StaticKeywordSplice)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q does not have splice ability", spliceCard.ID()))
		}
		can = false
		return can
	}
	subtypeString, ok := modifiers["Subtype"].(string)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q splice modifiers do not contain subtype", spliceCard.ID()))
		}
		can = false
		return can
	}
	subtype, ok := mtg.StringToSubtype(subtypeString)
	if !ok {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card %q has invalid subtype %q", spliceCard.ID(), subtypeString))
		}
		can = false
		return can
	}
	if !cardToCast.Match(has.Subtype(subtype)) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf(
				"card %q does not have subtype %q",
				cardToCast.ID(),
				subtype,
			))
		}
		can = false
	}
	return can
}
