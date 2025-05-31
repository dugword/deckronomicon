package judge

import (
	"deckronomicon/packages/engine/static_abilities"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"encoding/json"
	"fmt"
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
	object query.Object,
	ability gob.Ability,
	ruling *Ruling,
) bool {
	can := true
	if object.Controller() != player.ID() {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "permanent is not controlled by player")
		}
		can = false
	}
	// TODO: I don't like having to parse the cost here:
	c, err := cost.ParseCost(ability.Cost(), object)
	if err != nil {
		panic("failed to parse ability cost: " + err.Error())
		can = false // Skip abilities with invalid costs
	}
	if !CanPayCost(c, object, game, player) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot pay cost for ability: "+c.Description())
		}
		can = false
	}
	if ability.Speed() == mtg.SpeedSorcery {
		if !CanPlaySorcerySpeed(game, player.ID(), ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "ability cannot be activated at instant speed")
			}
			can = false
		}
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
	if !player.ZoneContains(zone, has.ID(card.ID())) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("player does not have card in zone %q", zone))
		}
		can = false
	}
	if !card.Match(is.Spell()) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "card is not a spell")
		}
		can = false
	}
	if !CanCastFromZone(game, player, zone, card, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("card cannot be cast from zone %q", zone))
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

func CanCastFromZone(
	game state.Game,
	player state.Player,
	zone mtg.Zone,
	card gob.Card,
	ruling *Ruling,
) bool {
	can := true
	switch zone {
	case mtg.ZoneHand:
		if !CanPayCost(card.ManaCost(), card, game, player) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "cannot pay cost for spell: "+card.ManaCost().Description())
			}
			can = false
		}
	case mtg.ZoneGraveyard:
		if !card.Match(has.StaticKeyword(mtg.StaticKeywordFlashback)) {
			can = false
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "card does not have flashback")
			}
			return can
		}
		var costString string
		// TODO: Maybe these should be methods on the Card type?
		// Or maybe a helper function?
		for _, ability := range card.StaticAbilities() {
			if ability.StaticKeyword() != mtg.StaticKeywordFlashback {
				continue
			}
			flashbackModifiers := static_abilities.FlashbackModifiers{}
			if err := json.Unmarshal(ability.Modifiers, &flashbackModifiers); err != nil {
				panic(fmt.Errorf("failed to unmarshal flashback modifiers: %w", err))
				can = false
			}
			costString = flashbackModifiers.Cost
		}
		if costString == "" {
			panic(fmt.Errorf("flashback ability on card %q is missing Cost modifier", card.Name()))
			can = false
		}
		costs, err := cost.ParseCost(costString, card)
		if err != nil {
			panic(fmt.Errorf("failed to parse cost: %w", err))
			can = false
		}
		if !CanPayCost(costs, card, game, player) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "cannot pay cost for spell: "+costs.Description())
			}
			can = false
		}
	}
	return can
}
