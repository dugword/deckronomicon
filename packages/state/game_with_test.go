package state

import (
	"deckronomicon/packages/game/definition"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWithResetPriorityPasses(t *testing.T) {
	got := Game{}.WithResetPriorityPasses()
	want := NewGameFromDefinition(definition.Game{
		PlayersPassedPriority: map[string]bool{},
	})
	if diff := cmp.Diff(want, got, AllowAllUnexported); diff != "" {
		t.Errorf("WithResetPriorityPasses() mismatch (-want +got):\n%s", diff)
	}
}

func TestWithPlayerPassedPriority(t *testing.T) {
	const playerID = "Test Player"
	got := Game{}.WithPlayerPassedPriority(playerID)
	want := NewGameFromDefinition(definition.Game{
		PlayersPassedPriority: map[string]bool{
			"Test Player": true,
		},
	})
	if diff := cmp.Diff(want, got, AllowAllUnexported); diff != "" {
		t.Errorf("WithPlayerPassedPriority(%s) mismatch (-want +got):\n%s", playerID, diff)
	}
}
