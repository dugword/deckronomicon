package actionparser

import (
	"deckronomicon/packages/engine/action"
	"fmt"
	"strconv"
	"strings"
)

func parseEmitMetric(
	metric string,
) (action.EmitMetricAction, error) {
	parts := strings.Fields(metric)
	if len(parts) != 2 {
		return action.EmitMetricAction{}, fmt.Errorf("emit metric command requires a metric name and a value")
	}
	metricName := parts[0]
	valueStr := parts[1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return action.EmitMetricAction{}, fmt.Errorf("emit value for metric %q must be numeric: %w", metricName, err)
	}

	return action.NewEmitMetricAction(metricName, value), nil
}
