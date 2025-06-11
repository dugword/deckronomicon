package engine

// TODO: Document what level things should live at. Maybe apply is where the
// core game engine logic and enforcement lives. it takes the structured
// imput, verifies per the rules of the game it can happen, and then applies
// it.

// TODO: Events should have small flat string values where possible because they get turned into JSON,
// we don't want to capture a lot of redundant information in the event.
// that means a lookkup to player will usually happen here.

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/mtg"
	"deckronomicon/packages/state"
	"errors"
	"fmt"
)

func (e *Engine) Apply(gameEvent event.GameEvent) error {
	e.log.Info("Applying event:", gameEvent.EventType())
	e.record.Add(gameEvent)
	game, err := e.applyEvent(e.game, gameEvent)
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
	return nil
}

func (e *Engine) applyEvent(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
	switch evnt := gameEvent.(type) {
	case event.GameLifecycleEvent:
		return e.applyGameLifecycleEvent(game, evnt)
	case event.GameStateChangeEvent:
		return e.applyGameStateChangeEvent(game, evnt)
	case event.PlayerEvent:
		return e.applyPlayerEvent(game, evnt)
	case event.PriorityEvent:
		return e.applyPriorityEvent(game, evnt)
	case event.TurnBasedActionEvent:
		return e.applyTurnBasedActionEvent(game, evnt)
	case event.MilestoneEvent:
		return game, nil
	case event.CheatEvent:
		return e.applyCheatEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown event type: %T", evnt)
	}
}
