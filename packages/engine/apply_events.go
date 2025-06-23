package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/reducer"
	"deckronomicon/packages/game/mtg"
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
	triggeredEvents, err := e.CheckTriggeredEffects(e.game, gameEvent)
	if err != nil {
		return fmt.Errorf("failed to check triggered effects: %w", err)
	}
	for _, triggeredEvent := range triggeredEvents {
		if err := e.ApplyEvent(triggeredEvent); err != nil {
			return fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
		}
	}
	return nil
}
