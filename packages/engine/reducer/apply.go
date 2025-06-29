package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"fmt"
	"slices"
)

func ApplyEventAndTriggers(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
	game, err := applyEvent(game, gameEvent)
	if err != nil {
		return game, fmt.Errorf("failed to apply event %q: %w", gameEvent.EventType(), err)
	}
	triggeredAbilities, err := CheckTriggeredAbilities(game, gameEvent)
	if err != nil {
		return game, fmt.Errorf("failed to check triggered abilities: %w", err)
	}
	for _, triggeredAbility := range triggeredAbilities {
		triggeredEvents, err := HandleTriggeredAbility(game, triggeredAbility, gameEvent)
		if err != nil {
			return game, fmt.Errorf("failed to handle triggered ability %q: %w", triggeredAbility.ID, err)
		}
		for _, triggeredEvent := range triggeredEvents {
			game, err = ApplyEventAndTriggers(game, triggeredEvent)
			if err != nil {
				return game, fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
			}
		}
	}
	// TODO: Make this more elegant and generic
	if gameEvent.EventType() == "BeginEndStep" {
		// Remove all triggered abilities that are set to end at the end of the turn.
		for _, ta := range game.TriggeredAbilities() {
			if ta.Duration == mtg.DurationEndOfTurn {
				game, err = ApplyEventAndTriggers(game, event.RemoveTriggeredAbilityEvent{
					ID: ta.ID,
				})
				if err != nil {
					return game, fmt.Errorf("failed to remove triggered ability %q: %w", ta.ID, err)
				}
			}
		}
	}
	return game, nil
}

func applyEvent(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
	switch evnt := gameEvent.(type) {
	case event.GameLifecycleEvent:
		return applyGameLifecycleEvent(game, evnt)
	case event.GameStateChangeEvent:
		return applyGameStateChangeEvent(game, evnt)
	case event.PlayerEvent:
		return applyPlayerEvent(game, evnt)
	case event.PriorityEvent:
		return applyPriorityEvent(game, evnt)
	case event.StackEvent:
		return applyStackEvent(game, evnt)
	case event.TriggeredAbilityEvent:
		return applyTriggeredAbilityEvent(game, evnt)
	case event.TurnBasedActionEvent:
		return applyTurnBasedActionEvent(game, evnt)
	case event.MilestoneEvent:
		return game, nil
	case event.CheatEvent:
		return applyCheatEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown event type: %T", evnt)
	}
}

func CheckTriggeredAbilities(game state.Game, evnt event.GameEvent) ([]gob.TriggeredAbility, error) {
	var triggeredAbilities []gob.TriggeredAbility
	for _, ta := range game.TriggeredAbilities() {
		if MatchesTrigger(ta.Trigger, evnt, game, ta.PlayerID) {
			triggeredAbilities = append(triggeredAbilities, ta)
		}
	}
	return triggeredAbilities, nil
}

func HandleTriggeredAbility(game state.Game, triggeredAbility gob.TriggeredAbility, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	var effectWithTargets []effect.EffectWithTarget
	for _, efct := range triggeredAbility.Effects {
		if addManaEffect, ok := efct.(effect.AddMana); ok {
			result, err := resolver.ResolveAddMana(game, triggeredAbility.PlayerID, addManaEffect)
			if err != nil {
				return nil, fmt.Errorf("failed to apply effect %q: %w", "AddMana", err)
			}
			events = append(events, result.Events...)
		} else {
			effectWithTargets = append(effectWithTargets, effect.EffectWithTarget{
				Effect: efct,
				Target: effect.Target{
					Type: mtg.TargetTypeNone,
				},
			})
		}
	}
	if len(effectWithTargets) > 0 {
		events = append(events, event.PutAbilityOnStackEvent{
			PlayerID:          triggeredAbility.PlayerID,
			SourceID:          triggeredAbility.SourceID,
			AbilityID:         triggeredAbility.ID,
			AbilityName:       "Triggered Effect",
			EffectWithTargets: effectWithTargets,
		})
	}
	return events, nil
}

func MatchesTrigger(trigger gob.Trigger, evnt event.GameEvent, game state.Game, playerID string) bool {
	// TODO: This match logic should live in the trigger itself I think, otherwise this is going to get out of hand.
	// Or maybe not because we have a generic "filter" in the trigger that is applied differently based on the event type.
	// Maybe this needs to be applied in a dispatching reducer pattern like the apply events function.
	// Maybe this should be in the judge package.
	// TODO: Yeah probably should be in the judge package.
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
		BeginEndStepEvent, ok := evnt.(event.BeginEndStepEvent)
		if !ok {
			return false
		}
		if BeginEndStepEvent.PlayerID != playerID {
			return false
		}
		return true
	}
	return false
}
