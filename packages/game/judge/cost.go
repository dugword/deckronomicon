package judge

import (
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/state"
)

func CanPayCost(someCost cost.Cost, game state.Game, player state.Player) bool {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		return canPayCompositeCost(c, game, player)
	case cost.ManaCost:
		return canPayManaCost(c, game, player)
	case cost.TapCost:
		return canPayTapCost(c, game, player)
	default:
		return false // Unsupported cost type
	}
}
func canPayCompositeCost(c cost.CompositeCost, game state.Game, player state.Player) bool {
	// Check if the player can pay all parts of the composite cost
	for _, subCost := range c.Costs() {
		if !CanPayCost(subCost, game, player) {
			return false
		}
	}
	return true
}

func canPayManaCost(c cost.ManaCost, game state.Game, player state.Player) bool {
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

func canPayTapCost(c cost.TapCost, game state.Game, player state.Player) bool {
	if !game.Battlefield().Contains(has.ID(c.Permanent().ID())) || c.Permanent().Controller() != player.ID() || c.Permanent().IsTapped() {
		return false
	}
	return true
}
