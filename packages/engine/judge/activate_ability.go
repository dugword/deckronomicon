package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
)

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
	if !CanPayCost(ability.Cost(), object, game, player, ruling) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot pay cost for ability: "+ability.Cost().Description())
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
