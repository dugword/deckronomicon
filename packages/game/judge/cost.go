package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
)

func CanPayCost(someCost cost.Cost, object query.Object, game state.Game, player state.Player) bool {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return canPayCompositeCost(c, object, game, player)
	case cost.ManaCost:
		return canPayManaCost(c, object, game, player)
	case cost.TapCost:
		return canPayTapCost(object, player, game)
	default:
		return false // Unsupported cost type
	}
}
func canPayCompositeCost(c cost.CompositeCost, object query.Object, game state.Game, player state.Player) bool {
	// Check if the player can pay all parts of the composite cost
	for _, subCost := range c.Costs() {
		if !CanPayCost(subCost, object, game, player) {
			return false
		}
	}
	return true
}

func canPayManaCost(c cost.ManaCost, object query.Object, game state.Game, player state.Player) bool {
	/*
		// Check if the player has enough mana to pay the cost
		availableMana := player.ManaAvailable()
		for _, mana := range cost.Mana() {
			if availableMana[mana.Color] < mana.Amount {
				return false
			}
		}
		return true
	*/
	return true
}

func canPayTapCost(object query.Object, player state.Player, game state.Game) bool {
	permanent, ok := object.(gob.Permanent)
	if !ok {
		return false // The object must be a permanent to pay a tap cost
	}
	if !game.Battlefield().Contains(has.ID(object.ID())) || object.Controller() != player.ID() || permanent.IsTapped() {
		return false
	}
	return true
}
