package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/judge"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/engine/resolver"
	"deckronomicon/packages/game/effect"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func (e *Engine) ApplyEvent(gameEvent event.GameEvent) error {
	e.log.Info("Applying event:", gameEvent.EventType())
	e.record.Add(gameEvent)
	game, err := reducer.ApplyEvent(e.game, gameEvent)
	if err != nil {
		if errors.Is(err, mtg.ErrGameOver) {
			e.log.Info("Game over detected")
			return err
		}
		e.log.Critical("Failed to apply event:", err)
		return fmt.Errorf("failed to apply event %q: %w", gameEvent.EventType(), err)
	}
	e.game = game
	for _, agent := range e.agents {
		agent.ReportState(e.game)
	}
	triggeredAbilities, err := e.CheckTriggeredAbilities(e.game, gameEvent)
	if err != nil {
		return fmt.Errorf("failed to check triggered abilities: %w", err)
	}
	for _, triggeredAbility := range triggeredAbilities {
		triggeredEvents, err := e.HandleTriggeredAbility(e.game, triggeredAbility, gameEvent)
		if err != nil {
			return fmt.Errorf("failed to handle triggered ability %q: %w", triggeredAbility.ID, err)
		}
		if triggeredAbility.Duration == mtg.DurationEndOfTurn && gameEvent.EventType() == "BeginEndStep" {
			triggeredEvents = append(triggeredEvents, event.RemoveTriggeredAbilityEvent{
				ID: triggeredAbility.ID,
			})
		}
		for _, triggeredEvent := range triggeredEvents {
			if err := e.ApplyEvent(triggeredEvent); err != nil {
				return fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
			}
		}
	}
	return nil
}

func (e *Engine) CheckTriggeredAbilities(game state.Game, evnt event.GameEvent) ([]gob.TriggeredAbility, error) {
	var triggeredAbilities []gob.TriggeredAbility
	for _, ta := range game.TriggeredAbilities() {
		if judge.MatchesTrigger(ta.Trigger, evnt, game, ta.PlayerID) {
			triggeredAbilities = append(triggeredAbilities, ta)
		}
	}
	return triggeredAbilities, nil
}

func (e *Engine) HandleTriggeredAbility(game state.Game, triggeredAbility gob.TriggeredAbility, evnt event.GameEvent) ([]event.GameEvent, error) {
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
