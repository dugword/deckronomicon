package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"fmt"
)

// canCastSpell checks for generic conditions that apply regardless of the zone.
func canCastSpell(
	game *state.Game,
	playerID string,
	zone mtg.Zone,
	card *gob.Card,
	cost cost.Cost,
	autoPayCost bool,
	autoPayColors []mana.Color,
	ruling *Ruling,
) bool {
	player := game.GetPlayer(playerID)
	can := true
	if !player.ZoneContains(zone, has.ID(card.ID())) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player does not have card in zone %q", zone))
		}
		can = false
	}
	if !card.Match(is.Not(is.Land())) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "card is not a spell")
		}
		can = false
	}
	if card.Match(query.Or(is.PermanentCardType(), has.CardType(mtg.CardTypeSorcery))) {
		if !CanPlaySorcerySpeed(game, playerID, ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "spell cannot be played at instant speed")
			}
			can = false
		}
	}
	if autoPayCost {
		if !CanPayCostAutomatically(cost, card, game, playerID, autoPayColors, ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "cannot pay cost for spell: "+cost.Description())
			}
			can = false
		}
	} else {
		if !CanPayCost(cost, card, game, playerID, ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "cannot pay cost for spell: "+cost.Description())
			}
			can = false
		}
	}
	return can
}

func CanCastSpellWithFlashback(
	game *state.Game,
	playerID string,
	card *gob.Card,
	cost cost.Cost,
	autoPayCost bool,
	autoPayColors []mana.Color,
	ruling *Ruling,
) bool {
	can := true
	if !card.Match(has.StaticKeyword(mtg.StaticKeywordFlashback)) {
		can = false
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "card does not have flashback")
		}
		return can
	}
	if !canCastSpell(game, playerID, mtg.ZoneGraveyard, card, cost, autoPayCost, autoPayColors, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot cast spell from graveyard")
		}
		return false
	}
	return can
}

func CanCastSpellFromHand(
	game *state.Game,
	playerID string,
	card *gob.Card,
	cost cost.Cost,
	autoPayCost bool,
	autoPayColors []mana.Color,
	ruling *Ruling,
) bool {
	can := true
	if !canCastSpell(game, playerID, mtg.ZoneHand, card, cost, autoPayCost, autoPayColors, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot cast spell from hand")
		}
		can = false
	}
	return can
}
