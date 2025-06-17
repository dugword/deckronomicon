package engine

import (
	"deckronomicon/packages/engine/effect"
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/target"
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

func (e *Engine) CheckTriggeredEffects(game state.Game, evnt event.GameEvent) ([]event.GameEvent, error) {
	var triggeredEvents []event.GameEvent
	for _, te := range game.TriggeredEffects() {
		fmt.Println("Checking triggered effect (sourceName):", te.SourceName, "for event:", evnt.EventType())
		if e.MatchesTrigger(te.Trigger, evnt, game, te.PlayerID) {
			events, err := e.HandleTriggeredEffect(game, te.PlayerID, te, evnt)
			if err != nil {
				return nil, err
			}
			triggeredEvents = append(triggeredEvents, events...)
		}
		// TODO Make this more generic and elegant.
		if te.Duration == mtg.DurationEndOfTurn && evnt.EventType() == "BeginEndStep" {
			triggeredEvents = append(triggeredEvents, event.RemoveTriggeredEffectEvent{
				ID: te.ID,
			})
		}
		fmt.Println("Removing triggered effect (sourceName):", te.SourceName, "for event:", evnt.EventType())
	}
	return triggeredEvents, nil
}

func (e *Engine) MatchesTrigger(trigger state.Trigger, evnt event.GameEvent, game state.Game, playerID string) bool {
	fmt.Println("Matches!")
	// TODO: This match logic should live in the trigger itself I think, otherwise this is going to get out of hand.
	// Or maybe not because we have a generic "filter" in the trigger that is applied differently based on the event type.
	// Maybe this needs to be applied in a dispatching reducer pattern like the apply events function.
	switch trigger.EventType {
	case "LandTappedForMana":
		LandTappedForManaEvent, ok := evnt.(event.LandTappedForManaEvent)
		if !ok {
			return false
		}
		if LandTappedForManaEvent.PlayerID != playerID {
			return false
		}
		if trigger.Filter.Subtypes != nil {
			for _, subtype := range trigger.Filter.Subtypes {
				if !slices.Contains(LandTappedForManaEvent.Subtypes, subtype) {
					return false
				}
			}
		}
		return true
	case "BeginEndStep":
		fmt.Println("Checking BeginEndStep trigger for player:", playerID)
		BeginEndStepEvent, ok := evnt.(event.BeginEndStepEvent)
		if !ok {
			return false
		}
		fmt.Println("Was ok")
		if BeginEndStepEvent.PlayerID != playerID {
			return false
		}
		fmt.Println("All checks passed")
		return true
	}
	return false
}

/*
func (e *Engine) HandleTriggeredEffectOld(game state.Game, playerID string, te state.TriggeredEffect, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	player, ok := game.GetPlayer(playerID)
	if !ok {
		e.log.Error("Player not found for ID:", playerID)
		return events, fmt.Errorf("player not found for ID: %s", playerID)
	}
	// TODO: Effects other than "AddMana" need to be put on the stack.
	for _, effectSpec := range te.Effect {
		effectEvents, err := e.ResolveEffect(game, player, nil, nil, effectSpec)
		if err != nil {
			e.log.Error("Failed to resolve effect:", effectSpec.Name, "Error:", err)
			continue
		}
		events = append(events, effectEvents...)
	}
	if te.OneShot {
		events = append(events, event.RemoveTriggeredEffectEvent{
			ID: te.ID,
		})
	}
	return events, nil
}
*/

func (e *Engine) HandleTriggeredEffect(game state.Game, playerID string, te state.TriggeredEffect, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	player, ok := game.GetPlayer(playerID)
	if !ok {
		e.log.Error("Player not found for ID:", playerID)
		return events, fmt.Errorf("player not found for ID: %s", playerID)
	}
	// TODO: This is redundante with the activate ability code and the resolve effect code.
	// We should refactor this to use the same logic.
	// I think what I am doing here though is the pattern to use for all abilities/effects.
	// I.e. resolving the "AddMana" effect and putting everything else as an ability on the stack.
	var effectSpecs []definition.EffectSpec
	for _, effectSpec := range te.Effect {
		if effectSpec.Name == "AddMana" {
			efct, err := effect.Build(effectSpec)
			if err != nil {
				return nil, fmt.Errorf("effect %q not found: %w", effectSpec.Name, err)
			}
			effectResults, err := efct.Resolve(game, player, nil, target.TargetValue{})
			if err != nil {
				return nil, fmt.Errorf("failed to apply effect %q: %w", effectSpec.Name, err)
			}
			events = append(events, effectResults.Events...)
		} else {
			effectSpecs = append(effectSpecs, effectSpec)
		}
	}
	if len(effectSpecs) > 0 {
		events = append(events, event.PutAbilityOnStackEvent{
			PlayerID:    playerID,
			SourceID:    te.ID,
			AbilityID:   te.ID,
			AbilityName: "Triggered Effect",
			EffectSpecs: effectSpecs,
		})
	}
	return events, nil
}
