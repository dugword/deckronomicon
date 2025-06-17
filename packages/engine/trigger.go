package engine

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
)

func (e *Engine) CheckTriggeredEffects(game state.Game, event event.GameEvent) []event.GameEvent {
	/*
		var triggeredEvents []event.GameEvent
		for _, te := range game.TriggeredEffects() {
			if e.MatchesTrigger(te.Trigger, event, game, player) {
				events := e.HandleTriggeredEffect(game, player, te, event)
				triggeredEvents = append(triggeredEvents, events...)
			}
		}
		return triggeredEvents
	*/
	return nil
}

func (e *Engine) MatchesTrigger(trigger state.Trigger, event event.GameEvent, game state.Game, player state.Player) bool {
	/*
		switch trigger.Type {
		case "TappedForMana":
			event, ok := event.(event.TappedForManaEvent)
			if !ok {
				return false
			}
			obj := game.GetObject(event.ObjectID)
			if trigger.Filter != nil {
				if !obj.Match(filter) {
					return false
				}
			}
		}
	*/
	return false
}

func (e *Engine) HandleTriggeredEffect(game state.Game, player state.Player, te state.TriggeredEffect, event event.GameEvent) []event.GameEvent {
	/*
		switch te.Effect.Type {
		case state.ResponseAddMana:
			mana := *te.Effect.Mana
			playerID := GetTriggereingPlayerID(triggeredEvents)
			return []state.GameEvent{
				state.AddManaEvent{
					PlayerID: playerID,
					Mana:     mana,
				}
			}

		}
	*/
	return nil
}
