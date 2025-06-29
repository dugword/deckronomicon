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

// AutoActivateManaSources attempts to automatically activate mana sources to pay for a given cost.
// It will iterate through the cost, activating mana sources for each color of mana required.
// The colors parameter specifies the order in which to try to pay the colors.
// It returns a slice of events that occurred during the process, or an error if it fails
// The events returned will include the acticvation of mana sources only.
// Spending events are not included in the events, as they are handled separately by the pay.Cost function.
func AutoActivateManaSources(game state.Game, someCost cost.Cost, object gob.Object, playerID string, colors []mana.Color) ([]event.GameEvent, error) {
	// The game state returned from withActiatedManaSources is not used, as the game state is updated
	// by the reducer.ApplyEvent function in the Engine's ApplyEvent method.
	// This interim game state is use to track the which mana sources will be activated
	// and the impact to the player's mana pool.
	switch c := someCost.(type) {
	// TODO: Assumes no nested composite costs
	case cost.CompositeCost:
		var events []event.GameEvent
		for _, subCost := range c.Costs() {
			switch sc := subCost.(type) {
			case cost.CompositeCost:
				panic("nested composite costs are not supported")
			case cost.ManaCost:
				var subEvents []event.GameEvent
				var err error
				game, subEvents, err = withActivateManaSources(game, playerID, sc, colors)
				if err != nil {
					return nil, err
				}
				events = append(events, subEvents...)
			}
		}
		return events, nil
	case cost.ManaCost:
		_, events, err := withActivateManaSources(game, playerID, c, colors)
		return events, err
	default:
		return nil, nil
	}
}

func withUpdatePlayerSpendAmount(
	game state.Game,
	playerID string,
	amount mana.Amount,
	colors []mana.Color,
) (state.Game, mana.Amount) {
	player := game.GetPlayer(playerID)
	manaPool := player.ManaPool()
	manaPool, remaining := manaPool.WithSpendAmount(amount, colors)
	player = player.WithManaPool(manaPool)
	game = game.WithUpdatedPlayer(player)
	return game, remaining
}

func withActivateManaSources(
	game state.Game,
	playerID string,
	manaCost cost.ManaCost,
	colors []mana.Color,
) (state.Game, []event.GameEvent, error) {
	var events []event.GameEvent
	var err error
	remaining := manaCost.Amount()
	// Pay down from existing mana pool
	game, remaining = withUpdatePlayerSpendAmount(
		game,
		playerID,
		remaining,
		colors,
	)
	// Pay for each color
	for _, c := range mana.Colors() {
		if remaining.AmountOf(c) <= 0 {
			continue
		}
		var activateEvents []event.GameEvent
		game, remaining, activateEvents, err = activateManaSourcesForColored(
			game,
			playerID,
			remaining,
			c,
		)
		if err != nil {
			return game, nil, err
		}
		events = append(events, activateEvents...)
	}
	// Pay for generic mana
	if remaining.Generic() > 0 {
		var activateEvents []event.GameEvent
		game, remaining, activateEvents, err = activateManaSourcesForGeneric(
			game,
			playerID,
			remaining,
			colors,
		)
		if err != nil {
			return game, nil, err
		}
		events = append(events, activateEvents...)
	}
	if remaining.Total() > 0 {
		return game, nil, errors.New("not enough mana in sources to auto-pay cost")
	}
	return game, events, nil
}

func activateManaSourcesForColored(
	game state.Game,
	playerID string,
	remaining mana.Amount,
	manaColor mana.Color,
) (state.Game, mana.Amount, []event.GameEvent, error) {
	var events []event.GameEvent
	var err error
	for remaining.AmountOf(manaColor) > 0 {
		lands := game.Battlefield().FindAll(
			query.And(is.Land(), is.Untapped(), has.Controller(playerID), has.HasManaAbility(manaColor)),
		)
		// Return when mana sources are exhausted
		if len(lands) == 0 {
			break
		}
		var activatedEvents []event.GameEvent
		game, activatedEvents, err = activateManaSource(
			game,
			playerID,
			lands[0].ID(),
		)
		events = append(events, activatedEvents...)
		game, remaining = withUpdatePlayerSpendAmount(
			game,
			playerID,
			remaining,
			[]mana.Color{manaColor},
		)
	}
	return game, remaining, events, err
}

func activateManaSourcesForGeneric(
	game state.Game,
	playerID string,
	remaining mana.Amount,
	colors []mana.Color,
) (state.Game, mana.Amount, []event.GameEvent, error) {
	var events []event.GameEvent
	var err error
	for _, manaColor := range colors {
		for remaining.Generic() > 0 {
			lands := game.Battlefield().FindAll(
				query.And(is.Land(), is.Untapped(), has.Controller(playerID), has.HasManaAbility(manaColor)),
			)
			// Check next color when mana sources are exhausted
			if len(lands) == 0 {
				break
			}
			var activatedEvents []event.GameEvent
			game, activatedEvents, err = activateManaSource(
				game,
				playerID,
				lands[0].ID(),
			)
			events = append(events, activatedEvents...)
			game, remaining = withUpdatePlayerSpendAmount(
				game,
				playerID,
				remaining,
				[]mana.Color{manaColor},
			)
		}
	}
	return game, remaining, events, err
}

// TODO: redundate with judge.GetAvailableMana
func activateManaSource(
	game state.Game,
	playerID string,
	landID string,
) (state.Game, []event.GameEvent, error) {
	player := game.GetPlayer(playerID)
	land, ok := game.Battlefield().Get(landID)
	if !ok {
		return game, nil, fmt.Errorf("land %s not found on battlefield", landID)
	}
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
		costEvents := Cost( // pay.Cost() of the mana ability
			ability.Cost(),
			land,
			player.ID(),
		)
		events = append(events, costEvents...)
		events = append(events, event.LandTappedForManaEvent{
			PlayerID: player.ID(),
			ObjectID: land.ID(),
			Subtypes: land.Subtypes(),
		})
		for _, event := range events {
			var err error
			game, err = reducer.ApplyEventAndTriggers(game, event)
			if err != nil {
				return game, nil, err
			}
		}
		for _, efct := range ability.Effects() {
			addManaEffect, ok := efct.(effect.AddMana)
			if !ok {
				continue
			}
			// TODO: Assumes no Resume function or choices required for the effect.
			// This logic should be handled centrally as there are a few places where
			// mana abilities are resolved and this logic is duplicated.
			result, err := resolver.ResolveAddMana(game, player.ID(), addManaEffect)
			if err != nil {
				return game, nil, fmt.Errorf("failed to resolve add mana effect: %w", err)
			}
			if result.Resume != nil {
				return game, nil, fmt.Errorf("mana ability %s requires choices to be made", ability.Name())
			}
			events = append(events, result.Events...)
			for _, event := range result.Events {
				game, err = reducer.ApplyEventAndTriggers(game, event)
				if err != nil {
					return game, nil, err
				}
			}
			return game, events, nil
		}
	}
	return game, nil, nil
}
