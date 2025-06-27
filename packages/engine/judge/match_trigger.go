package judge

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/game/gob"
	"deckronomicon/packages/state"
	"slices"
)

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
