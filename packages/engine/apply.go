package engine

import (
	"deckronomicon/packages/engine/event"
	"fmt"
)

func (e *Engine) Apply(event event.GameEvent) {
	e.record.AddEvent(event)
	e.state = ApplyEvent(e.state, event)
}

func ApplyEvent(state *GameState, event event.GameEvent) *GameState {
	switch event.Type {
	case "PassPriority":
		for _, player := range state.players {
			if player.id == event.PlayerID {
				player.passed = true
			}
		}
	}
	return state
}

func ApplyDrawCardEvent(
	state GameState,
	event event.GameEvent,
) (GameState, event.GameEvent, error) {
	player, err := state.GetPlayer(event.PlayerID)
	if err != nil {
		return state, event, fmt.Errorf("draw: %w", err)
	}
	if len(player.library) == 0 {
		return state, event, fmt.Errorf("draw: player %s has no cards to draw", player.id)
	}
	// Pop the top card from the library
	drawnCard := player.library[0]
	player.library = player.library[1:]
	// Add it to their hand
	player.hand = append(player.hand, drawnCard)
	// Update the player in game state
	state = state.WithUpdatedPlayer(player)
	// Attach the card ID to the GameEvent metadata for auditing/logging
	/*
		if event.Metadata == nil {
			event.Metadata = make(map[string]any)
		}
	*/
	/// event.Metadata["card_id"] = drawnCard.ID
	return state, event, nil
}
