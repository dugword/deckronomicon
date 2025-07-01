package engine

import (
	"deckronomicon/packages/engine/event"
	"encoding/json"
	"fmt"
)

type GameRecord struct {
	Seed   int64
	Events []event.GameEvent
	// FinalState GameStateSnapshot
}

func NewGameRecord(seed int64) *GameRecord {
	return &GameRecord{
		Seed:   seed,
		Events: []event.GameEvent{},
	}
}

func (r *GameRecord) Add(e event.GameEvent) {
	r.Events = append(r.Events, e)
}

type ExportEvent struct {
	EventType string
	Details   event.GameEvent
}

func (e ExportEvent) MarshalJSON() ([]byte, error) {
	detailsBytes, err := json.Marshal(e.Details)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(detailsBytes, &m); err != nil {
		return nil, err
	}
	if m["EventType"] != nil {
		return nil, fmt.Errorf("field 'EventType' already exists in event details")
	}
	m["EventType"] = e.EventType
	return json.Marshal(m)
}

func (r *GameRecord) Export() []ExportEvent {
	var events []ExportEvent
	for _, e := range r.Events {
		events = append(events, ExportEvent{
			EventType: e.EventType(),
			Details:   e,
		})
	}
	return events
}

func (r *GameRecord) ExportAnalytics() []ExportEvent {
	var events []ExportEvent
	for _, e := range r.Events {
		if _, ok := e.(event.AnalyticsEvent); ok {
			events = append(events, ExportEvent{
				EventType: e.EventType(),
				Details:   e,
			})
		}
	}
	return events
}
