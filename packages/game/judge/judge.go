package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"strings"
)

type Ruling struct {
	Explain bool
	Reasons []string
}

func (r *Ruling) Why() string {
	if r == nil || !r.Explain {
		return ""
	}
	if len(r.Reasons) == 0 {
		return "No reasons provided."
	}
	return "reasons: " + strings.Join(r.Reasons, ", ")
}

func CanActivateAbility(
	game state.Game,
	player state.Player,
	permanent gob.Permanent,
	ability gob.Ability,
	ruling *Ruling,
) bool {
	can := true
	if permanent.Controller() != player.ID() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "permanent is not controlled by player")
		}
		can = false
	}
	// TODO: I don't like having to parse the cost here:
	c, err := cost.ParseCost(ability.Cost(), permanent)
	if err != nil {
		panic("failed to parse ability cost: " + err.Error())
		can = false // Skip abilities with invalid costs
	}
	if !CanPayCost(c, game, player) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot pay cost for ability: "+c.Description())
		}
		can = false
	}
	return can
}

// TODO should this live on engine?
func CanPlaySorcerySpeed(game state.Game, playerID string, ruling *Ruling) bool {
	can := true
	if !game.IsStackEmtpy() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "stack is not empty")
		}
		can = false
	}
	if game.ActivePlayerID() != playerID {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "player is not active player")
		}
		can = false
	}
	if !(game.Step() == mtg.StepPrecombatMain ||
		game.Step() == mtg.StepPostcombatMain) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "not in a main phase")
		}
		can = false
	}
	return can
}

// TODO: Should this live on engine?
func CanCastSpell(
	game state.Game,
	player state.Player,
	zone mtg.Zone,
	card gob.Card,
	ruling *Ruling,
) bool {
	can := true
	if !card.Match(is.Spell()) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "card is not a spell")
		}
		can = false
	}
	if zone != mtg.ZoneHand {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "zone is not hand")
		}
		can = false
	}
	if card.Match(query.Or(is.Permanent(), has.CardType(mtg.CardTypeSorcery))) {
		if !CanPlaySorcerySpeed(game, player.ID(), ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "spell cannot be played at instant speed")
			}
			can = false
		}
	}
	return can
}

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
