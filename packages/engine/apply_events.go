package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/game/target"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
	"slices"
)

func (e *Engine) ApplyEvent(gameEvent event.GameEvent) error {
	e.log.Debug("Applying event:", gameEvent.EventType())
	game, record, err := innerApplyEvent(e.game, gameEvent)
	if err != nil {
		if errors.Is(err, mtg.ErrGameOver) {
			e.log.Debug("Game over detected")
			return err
		}
		e.log.Critical("Failed to apply event:", err)
		return err
	}
	for _, recordEvent := range record {
		e.record.Add(recordEvent)
	}
	e.game = game
	for _, agent := range e.agents {
		agent.ReportState(e.game)
	}
	return nil
}

// TODO: Would it be useful to return the record to see what was applied?
// func (e *Engine) MaybeApplyEvent(gameEvent event.GameEvent) (*state.Game, []event.GameEvent, error) {
func MaybeApplyEvent(game *state.Game, gameEvent event.GameEvent) (*state.Game, error) {
	game, _, err := innerApplyEvent(game, gameEvent)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func innerApplyEvent(game *state.Game, gameEvent event.GameEvent) (*state.Game, []event.GameEvent, error) {
	var record []event.GameEvent
	// Apply Events
	game, err := reducer.Apply(game, gameEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply event %q: %w", gameEvent.EventType(), err)
	}

	// Check Triggered Abilities
	triggeredAbilities, err := CheckTriggeredAbilities(game, gameEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check triggered abilities: %w", err)
	}
	for _, triggeredAbility := range triggeredAbilities {
		triggeredEvents, err := HandleTriggeredAbility(game, triggeredAbility, gameEvent)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle triggered ability %q: %w", triggeredAbility.ID, err)
		}
		// TODO: Check for infinite loops here
		for _, triggeredEvent := range triggeredEvents {
			gameTriggered, recordTriggered, err := innerApplyEvent(game, triggeredEvent)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
			}
			game = gameTriggered
			record = append(record, recordTriggered...)
		}
	}

	// TODO: Make this more elegant and generic
	if gameEvent.EventType() == "BeginEndStep" {
		// Remove all triggered abilities that are set to end at the end of the turn.
		for _, ta := range game.TriggeredAbilities() {
			if ta.Duration == mtg.DurationEndOfTurn {
				gameRemoveTriggered, recordRemoveTriggered, err := innerApplyEvent(game, &event.RemoveTriggeredAbilityEvent{
					ID: ta.ID,
				})
				if err != nil {
					return nil, nil, fmt.Errorf("failed to remove triggered ability %q: %w", ta.ID, err)
				}
				game = gameRemoveTriggered
				record = append(record, recordRemoveTriggered...)
			}
		}
	}
	return game, record, nil
}

func CheckTriggeredAbilities(game *state.Game, evnt event.GameEvent) ([]gob.TriggeredAbility, error) {
	var triggeredAbilities []gob.TriggeredAbility
	for _, ta := range game.TriggeredAbilities() {
		if MatchesTrigger(ta.Trigger, evnt, game, ta.PlayerID) {
			triggeredAbilities = append(triggeredAbilities, ta)
		}
	}
	return triggeredAbilities, nil
}

func HandleTriggeredAbility(game *state.Game, triggeredAbility gob.TriggeredAbility, evnt event.GameEvent) ([]event.GameEvent, error) {
	var events []event.GameEvent
	var effectWithTargets []*effect.EffectWithTarget
	for _, efct := range triggeredAbility.Effects {
		if addManaEffect, ok := efct.(*effect.AddMana); ok {
			result, err := resolver.ResolveAddMana(game, triggeredAbility.PlayerID, addManaEffect)
			if err != nil {
				return nil, fmt.Errorf("failed to apply effect %q: %w", "AddMana", err)
			}
			events = append(events, result.Events...)
		} else {
			effectWithTargets = append(effectWithTargets, &effect.EffectWithTarget{
				Effect: efct,
				Target: target.Target{
					Type: mtg.TargetTypeNone,
				},
			})
		}
	}
	if len(effectWithTargets) > 0 {
		events = append(events, &event.PutAbilityOnStackEvent{
			PlayerID:          triggeredAbility.PlayerID,
			SourceID:          triggeredAbility.SourceID,
			AbilityID:         triggeredAbility.ID,
			AbilityName:       "Triggered Effect",
			EffectWithTargets: effectWithTargets,
		})
	}
	return events, nil
}

func MatchesTrigger(trigger gob.Trigger, evnt event.GameEvent, game *state.Game, playerID string) bool {
	// TODO: This match logic should live in the trigger itself I think, otherwise this is going to get out of hand.
	// Or maybe not because we have a generic "filter" in the trigger that is applied differently based on the event type.
	// Maybe this needs to be applied in a dispatching reducer pattern like the apply events function.
	// Maybe this should be in the judge package.
	// TODO: Yeah probably should be in the judge package.
	switch trigger.EventType {
	case "LandTappedForMana":
		LandTappedForManaEvent, ok := evnt.(*event.LandTappedForManaEvent)
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
		BeginEndStepEvent, ok := evnt.(*event.BeginEndStepEvent)
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
