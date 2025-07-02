package action

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/engine/resenv"
	"deckronomicon/packages/state"
)

type EmitMetricAction struct {
	metric string
	value  int
}

func NewEmitMetricAction(metric string, value int) EmitMetricAction {
	return EmitMetricAction{
		metric: metric,
		value:  value,
	}
}

func (a EmitMetricAction) Name() string {
	return "Emit a metric"
}

func (a EmitMetricAction) Complete(game *state.Game, playerID string, resEnv *resenv.ResEnv) ([]event.GameEvent, error) {
	return []event.GameEvent{
		&event.EmitMetricEvent{
			PlayerID: playerID,
			Metric:   a.metric,
			Value:    a.value,
		},
	}, nil
}
