package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

func Apply(game *state.Game, gameEvent event.GameEvent) (*state.Game, error) {
	switch evnt := gameEvent.(type) {
	case event.AnalyticsEvent:
		return applyAnalyticsEvent(game, evnt)
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
	case event.TriggeredAbilityEvent:
		return applyTriggeredAbilityEvent(game, evnt)
	case event.TurnBasedActionEvent:
		return applyTurnBasedActionEvent(game, evnt)
	case *event.MilestoneEvent:
		return game, nil
	case event.CheatEvent:
		return applyCheatEvent(game, evnt)
	default:
		return game, fmt.Errorf("unknown event type: %T", evnt)
	}
}
