package judge

import (
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/store"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/state"
)

// TODO: Maybe check for some errors here and pass those through?
// E.g. is it an error because the player doesn't have enough mana,
// or is it an error because of some broken game state?
// TODO: Pass in ruling here and log which costs could not be paid
func CanPayCost(
	someCost cost.Cost,
	object gob.Object,
	game *state.Game,
	playerID string,
	ruling *Ruling,
) bool {
	costEvents, err := pay.Cost(someCost, object, playerID)
	if err != nil {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "unable to pay cost: "+err.Error())
		}
		return false
	}
	canPay := true
	str := store.NewStore(game, nil)
	for _, costEvent := range costEvents {
		if err := str.Apply(costEvent); err != nil {
			if ruling != nil && ruling.Explain {
				// TODO: Find some way to tie this to the specific cost event
				// Maybe each event needs to have a reason for why it's being applied.
				// E.g. the LoseLifeEvent could have a reason like "paying cost for spell X"
				ruling.Reasons = append(ruling.Reasons, "unable to pay cost requiring "+costEvent.EventType())
			}
			canPay = false
			break
		}
	}
	return canPay
}

// TODO: I think this is only checking for mana costs not all costs
func CanPayCostAutomatically(
	someCost cost.Cost,
	object gob.Object,
	game *state.Game,
	playerID string,
	colors []mana.Color,
	ruling *Ruling,
) bool {
	canPay := true
	// TODO: events no longer includes paying costs, this is just the activation of mana sources
	// still need to call the pay function to get the events
	events, err := pay.AutoActivateManaSources(game, someCost, object, playerID, colors)
	if err != nil {
		if ruling != nil && ruling.Explain {
			ruling.Reasons = append(ruling.Reasons, "unable to pay cost automatically: "+err.Error())
		}
		canPay = false
		return canPay
	}
	// Why am I applying these events?
	// Don't I apply them in auto activate?
	for _, event := range events {
		str := store.NewStore(game, nil)
		if err := str.Apply(event); err != nil {
			if ruling != nil && ruling.Explain {
				ruling.Reasons = append(ruling.Reasons, "unable to apply event for cost: "+event.EventType())
			}
			canPay = false
		}
	}
	return canPay
}

/*
func CanPotentiallyPayCost(someCost cost.Cost, object gob.Object, game *state.Game, playerID string, ruling *Ruling) bool {
	player := game.GetPlayer(playerID)
	potentialManaPool := GetAvailableMana(game, playerID)
	player = player.WithManaPool(potentialManaPool)
	game = game.WithUpdatedPlayer(player)
	return CanPayCost(someCost, object, game, playerID, ruling)
}
*/
