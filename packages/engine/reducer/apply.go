package reducer

// TODO: Document what level things should live at. Maybe apply is where the
// core game engine logic and enforcement lives. it takes the structured
// imput, verifies per the rules of the game it can happen, and then applies
// it.

// TODO: Events should have small flat string values where possible because they get turned into JSON,
// we don't want to capture a lot of redundant information in the event.
// that means a lookkup to player will usually happen here.

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

func ApplyEvent(game state.Game, gameEvent event.GameEvent) (state.Game, error) {
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
