package judge

import (
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
)

// TODO: For function signatures think about when to use little > big vs big > little
// e.g. game, playerID, zone, card vs card, zone playerID, game

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetLandsAvailableToPlay(game state.Game, player state.Player) []gob.CardInZone {
	var availableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if CanPlayLand(game, player, mtg.ZoneHand, card, nil) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if CanPlayLand(game, player, mtg.ZoneGraveyard, card, nil) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

func GetSpellsAvailableToCast(game state.Game, player state.Player) []gob.CardInZone {
	var availableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if CanCastSpell(game, player, mtg.ZoneHand, card, nil) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if CanCastSpell(game, player, mtg.ZoneGraveyard, card, nil) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetAbilitiesAvailableToActivate(game state.Game, player state.Player, ruling *Ruling) []gob.AbilityInZone {
	var availableAbilities []gob.AbilityInZone
	for _, permanent := range game.Battlefield().GetAll() {
		for _, ability := range permanent.ActivatedAbilities() {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, fmt.Sprintf("[ability %q]: ", ability.Name()))
			}
			if CanActivateAbility(game, player, permanent, ability, ruling) {
				availableAbilities = append(availableAbilities, gob.NewAbilityInZone(ability, permanent, mtg.ZoneBattlefield))
			}
		}
	}
	return availableAbilities
}
