package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
)

// TODO: Maybe these should be functions instead of methods on Game?

func CanActivateAbility(game state.Game, playerID string, permanent gob.Permanent, ability gob.Ability) bool {
	return true
}

// TODO should this live on engine?
func CanCastSorcery(game state.Game, playerID string) bool {
	if !game.IsStackEmtpy() {
		return false
	}
	if game.ActivePlayerID() != playerID {
		return false
	}
	if !(game.Step() == mtg.StepPrecombatMain ||
		game.Step() == mtg.StepPostcombatMain) {
		return false
	}
	return true
}

// TODO: Should this live on engine?
func CanPlayCard(
	game state.Game,
	player state.Player,
	zone mtg.Zone,
	card gob.Card,
) bool {
	if zone == mtg.ZoneGraveyard {
		return false
	}
	if card.Match(query.Or(is.Permanent(), has.CardType(mtg.CardTypeSorcery))) {
		if !CanCastSorcery(game, player.ID()) {
			return false
		}
		if card.Match(is.Land()) && player.LandPlayedThisTurn() {
			return false
		}
	}
	return true
}
