package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
)

// TODO: For function signatures think about when to use little > big vs big > little
// e.g. game, playerID, zone, card vs card, zone playerID, game

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetCardsAvailableToPlay(game state.Game, player state.Player) []gob.CardInZone {
	var availableCards []gob.CardInZone
	for _, card := range player.Hand().GetAll() {
		if CanPlayCard(game, player, mtg.ZoneHand, card) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneHand))
		}
	}
	for _, card := range player.Graveyard().GetAll() {
		if CanPlayCard(game, player, mtg.ZoneGraveyard, card) {
			availableCards = append(availableCards, gob.NewCardInZone(card, mtg.ZoneGraveyard))
		}
	}
	return availableCards
}

// TODO: This probably shouldn't be in Judge, maybe judge is just boolean functions
// that return true or false, and this should be where it's used, like in the agent functions
func GetAbilitiesAvailableToActivate(game state.Game, player state.Player) []gob.AbilityInZone {
	var availableAbilities []gob.AbilityInZone
	for _, permanent := range game.Battlefield().GetAll() {
		for _, ability := range permanent.ActivatedAbilities() {
			if permanent.Controller() != player.ID() {
				continue
			}
			// TODO: I don't like having to parse the cost here:
			c, err := cost.ParseCost(ability.Cost(), permanent)
			if err != nil {
				panic("failed to parse ability cost: " + err.Error())
				continue // Skip abilities with invalid costs
			}
			if !CanPayCost(c, game, player) {
				continue
			}
			if CanActivateAbility(game, player.ID(), permanent, ability) {
				availableAbilities = append(availableAbilities, gob.NewAbilityInZone(ability, permanent, mtg.ZoneBattlefield))
			}
		}
	}
	return availableAbilities
}
