package engine

import "deckronomicon/packages/engine/event"

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
