package judge

import (
	"deckronomicon/packages/engine/pay"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/query"
	"deckronomicon/packages/state"
)

// TODO: Maybe check for some errors here and pass those through?
// E.g. is it an error because the player doesn't have enough mana,
// or is it an error because of some broken game state?
// TODO: Pass in ruling here and log which costs could not be paid
func CanPayCost(someCost cost.Cost, object query.Object, game state.Game, player state.Player, ruling *Ruling) bool {
	costEvents := pay.PayCost(someCost, object, player)
	canPay := true
	for _, costEvent := range costEvents {
		var err error
		game, err = reducer.ApplyEvent(game, costEvent)
		if err != nil {
			if ruling != nil && ruling.Explain {
				// TODO: Find some way to tie this to the specific cost event
				// Maybe each event needs to have a reason for why it's being applied.
				// E.g. the LoseLifeEvent could have a reason like "paying cost for spell X"
				ruling.Reasons = append(ruling.Reasons, "unable to pay cost requiring "+costEvent.EventType())
			}
			canPay = false
		}
	}
	return canPay
}
