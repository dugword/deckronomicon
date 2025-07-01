package event

const (
	EventTypeLogMessage = "LogMessage"
	EventTypeEmitMetric = "EmitMetric"
)

type AnalyticsEvent interface{ isAnalyticsEvent() }

type AnalyticsBaseEvent struct{}

func (e AnalyticsBaseEvent) isAnalyticsEvent() {}

type LogMessageEvent struct {
	AnalyticsBaseEvent
	PlayerID string
	Message  string
}

func (e LogMessageEvent) EventType() string {
	return EventTypeLogMessage
}

type EmitMetricEvent struct {
	AnalyticsBaseEvent
	PlayerID string
	Metric   string
	Value    int
}

func (e EmitMetricEvent) EventType() string {
	return EventTypeEmitMetric
}
