package reducer

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"fmt"
)

func applyAnalyticsEvent(game *state.Game, evnt event.AnalyticsEvent) (*state.Game, error) {
	switch evnt.(type) {
	case *event.LogMessageEvent:
		return game, nil
	case *event.EmitMetricEvent:
		return game, nil
	default:
		return game, fmt.Errorf("unknown triggered event type '%T'", evnt)
	}
}
