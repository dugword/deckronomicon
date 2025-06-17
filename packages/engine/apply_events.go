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
		if err := agent.ReportState(e.game); err != nil {
			e.log.Critical("Failed to report state to agent:", err)
			return fmt.Errorf("failed to report state to agent for %q: %w", agent.PlayerID(), err)
		}
	}
	triggeredEvents := e.CheckTriggeredEffects(e.game, gameEvent)
	for _, triggeredEvent := range triggeredEvents {
		if err := e.ApplyEvent(triggeredEvent); err != nil {
			return fmt.Errorf("failed to apply triggered event %q: %w", triggeredEvent.EventType(), err)
		}
	}
	return nil
}
