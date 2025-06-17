package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

func (e *Engine) CheckTriggeredEffects(game state.Game, evnt event.GameEvent) ([]event.GameEvent, error) {
	var triggeredEvents []event.GameEvent
	for _, te := range game.TriggeredEffects() {
		if e.MatchesTrigger(te.Trigger, evnt, game, te.PlayerID) {
			events, err := e.HandleTriggeredEffect(game, te.PlayerID, te, evnt)
			if err != nil {
				return nil, err
			}
			triggeredEvents = append(triggeredEvents, events...)
		}
	}
	return triggeredEvents, nil
}

func (e *Engine) MatchesTrigger(trigger state.Trigger, evnt event.GameEvent, game state.Game, playerID string) bool {
	switch trigger.EventType {
	case "LandTappedForMana":
		fmt.Println("Checking LandTappedForMana trigger")
		LandTappedForManaEvent, ok := evnt.(event.LandTappedForManaEvent)
		if !ok {
			return false
		}
		fmt.Printf("Checking if event PlayerID %q matches trigger PlayerID %q\n", LandTappedForManaEvent.PlayerID, playerID)
		if LandTappedForManaEvent.PlayerID != playerID {
			return false
		}
		fmt.Println("Checking Subtypes")
		if trigger.Filter.Subtypes != nil {
			for _, subtype := range trigger.Filter.Subtypes {
				fmt.Println("Checking subtype", subtype)
				fmt.Println("LandTappedForManaEvent Subtypes", LandTappedForManaEvent.Subtypes)
				if !slices.Contains(LandTappedForManaEvent.Subtypes, subtype) {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (e *Engine) HandleTriggeredEffect(game state.Game, playerID string, te state.TriggeredEffect, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	player, ok := game.GetPlayer(playerID)
	if !ok {
		e.log.Error("Player not found for ID:", playerID)
		return events, fmt.Errorf("player not found for ID: %s", playerID)
	}
	for _, effectSpec := range te.Effect {
		effectEvents, err := e.ResolveEffect(game, player, nil, nil, effectSpec)
		if err != nil {
			e.log.Error("Failed to resolve effect:", effectSpec.Name, "Error:", err)
			continue
		}
		events = append(events, effectEvents...)
	}
	return events, nil
}
