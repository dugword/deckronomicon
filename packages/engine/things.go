package engine

import "deckronomicon/packages/engine/event"

func (e *Engine) PassPriority(playerID string) {
	// Log the action of passing priority
	// e.logger.Debug("Player %s passed priority", playerID)
	// Create a game event for passing priority
	event := event.NewPassPriorityEvent(playerID)
	// Apply the event to the game state
	e.Apply(event)
}
