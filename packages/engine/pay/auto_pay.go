package pay

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/cost"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mana"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/query"
	"deckronomicon/packages/query/has"
	"deckronomicon/packages/query/is"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

// WARNING: THIS IS BROKEN AND WILL NOT WORK AS EXPECTED
// This currently only accounts for mana costs that are paid by tapping lands.
// It does not handle other costs like sacrificing permanents, paying life, etc.

// TODO: This is really rough and needs to be cleaned up.
// Should optimtize to be smart about which lands to tap first,
// e.g. it should prefer single-color lands over dual-color lands,

// AutoPay attempts to automatically pay the given cost using the player's mana pool.
// It will tap lands for mana until the cost is fully paid or no more lands can be tapped.
// It returns a slice of events that occurred during the process, or an error if it fails.
//
// It will try to pay for each color of mana in the cost, and if it cannot find enough mana, it will return an error.
// If it successfully pays the cost, it returns the events generated during the process.

// TODO: Come up with a strategy for when to pass player and when to pass playerID
func AutoPay(game state.Game, playerID string, object gob.Object, cost cost.Cost) ([]event.GameEvent, error) {
	var events []event.GameEvent
	var err error
	player := game.GetPlayer(playerID)
	manaPool := player.ManaPool()
	manaCosts := GetManaCosts(cost)
	totalManaCost := GetTotalFromCosts(manaCosts)
	for totalManaCost.Total() > 0 {
		manaPool, totalManaCost = manaPool.WithPayDownFromAmount(totalManaCost, mana.Colors())
		player = player.WithManaPool(manaPool)
		game = game.WithUpdatedPlayer(player)
		if totalManaCost.White() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.White)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if totalManaCost.Blue() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.Blue)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if totalManaCost.Black() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.Black)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if totalManaCost.Red() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.Red)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if totalManaCost.Green() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.Green)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if totalManaCost.Colorless() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForColorAndSeeWhatHappens(game, playerID, mana.Colorless)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if len(events) == 0 {
			break // No more lands to tap or no mana abilities available
		}
	}
	for totalManaCost.Total() > 0 {
		manaPool, totalManaCost = manaPool.WithPayDownFromAmount(totalManaCost, mana.Colors())
		player = player.WithManaPool(manaPool)
		game = game.WithUpdatedPlayer(player)
		if totalManaCost.Generic() > 0 {
			var moreEvents []event.GameEvent
			game, moreEvents, err = TapALandForGenericAndSeeWhatHappens(game, playerID)
			if err != nil {
				return nil, err
			}
			events = append(events, moreEvents...)
			if len(moreEvents) > 0 {
				continue
			}
		}
		if len(events) == 0 {
			break // No more lands to tap or no mana abilities available
		}
	}
	if totalManaCost.Total() > 0 {
		return nil, errors.New("not enough mana to auto-pay cost")
	}
	return events, nil
}

func TapALandForColorAndSeeWhatHappens(
	game state.Game,
	playerID string,
	manaColor mana.Color,
) (state.Game, []event.GameEvent, error) {
	lands := game.Battlefield().FindAll(
		query.And(is.Land(), is.Untapped(), has.Controller(playerID), has.HasManaAbility(manaColor)),
	)
	return TapLandsAndSeeWhatHappens(
		game,
		playerID,
		lands,
	)
}

func TapALandForGenericAndSeeWhatHappens(
	game state.Game,
	playerID string,
) (state.Game, []event.GameEvent, error) {
	lands := game.Battlefield().FindAll(
		query.And(is.Land(), is.Untapped(), has.Controller(playerID)),
	)
	return TapLandsAndSeeWhatHappens(
		game,
		playerID,
		lands,
	)
}

func TapLandsAndSeeWhatHappens(
	game state.Game,
	playerID string,
	lands []gob.Permanent,
) (state.Game, []event.GameEvent, error) {
	player := game.GetPlayer(playerID)
	for _, land := range lands {
		for _, ability := range land.ActivatedAbilities() {
			if !ability.Match(is.ManaAbility()) {
				continue
			}
			events := []event.GameEvent{
				event.ActivateAbilityEvent{
					PlayerID:  player.ID(),
					SourceID:  land.ID(),
					AbilityID: ability.Name(),
					Zone:      mtg.ZoneBattlefield,
				},
			}
			for _, efct := range ability.Effects() {
				addManaEffect, ok := efct.(effect.AddMana)
				if !ok {
					continue
				}
				costEvents := Cost(ability.Cost(), land, player)
				events = append(events, costEvents...)
				events = append(events, event.LandTappedForManaEvent{
					PlayerID: player.ID(),
					ObjectID: land.ID(),
					Subtypes: land.Subtypes(),
				})
				result, err := resolver.ResolveAddMana(game, player.ID(), addManaEffect)
				if err != nil {
					return game, nil, fmt.Errorf("failed to resolve add mana effect: %w", err)
				}
				events = append(events, result.Events...)
				for _, event := range events {
					game, err = reducer.ApplyEvent(game, event)
					if err != nil {
						return game, nil, err
					}
				}
				return game, events, nil
			}
		}
	}
	return game, nil, nil // No valid land found or no mana ability
}

func GetTotalFromCosts(someCost []cost.ManaCost) mana.Amount {
	total := mana.Amount{}
	for _, c := range someCost {
		total = total.WithAddedAmount(c.Amount())
	}
	return total
}

func GetManaCosts(someCost cost.Cost) []cost.ManaCost {
	switch c := someCost.(type) {
	case cost.CompositeCost:
		var manaCosts []cost.ManaCost
		for _, subCost := range c.Costs() {
			manaCosts = append(manaCosts, GetManaCosts(subCost)...)
		}
		return manaCosts
	case cost.ManaCost:
		return []cost.ManaCost{c}
	default:
		return nil // No mana costs found
	}
}
