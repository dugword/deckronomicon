package judge

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

func CanActivateAbility(
	game *state.Game,
	playerID string,
	object gob.Object,
	ability *gob.Ability,
	ruling *Ruling,
	maybeApply func(game *state.Game, event event.GameEvent) (*state.Game, error),
) bool {
	can := true
	if object.Controller() != playerID {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "permanent is not controlled by player")
		}
		can = false
	}
	if !CanPayCost(ability.Cost(), object, game, playerID, ruling, maybeApply) {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "cannot pay cost for ability: "+ability.Cost().Description())
		}
		can = false
	}
	if ability.Speed() == mtg.SpeedSorcery {
		if !CanPlaySorcerySpeed(game, playerID, ruling) {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "ability cannot be activated at instant speed")
			}
			can = false
		}
	}
	return can
}
