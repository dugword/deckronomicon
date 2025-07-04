package judge

import (
	"deckronomicon/packages/engine/event"
	"deckronomicon/packages/state"
	"testing"
)

var mockApply = func(game *state.Game, event event.GameEvent) (*state.Game, error) {
	return game, nil
}

func TestGetAvailableMana(t *testing.T) {
	// TODO: This tests relies on the game state being updated correctly,
	// by the reducer. Which is not what we should be testing here.
	/*
		playerID := "Test Player"
		forestDefinition := definitiontest.ForestDefinition("Forest ID", playerID)
		islandDefinition := definitiontest.IslandDefinition("Island ID", playerID)
		islandDefinition.Tapped = true
		game := state.NewGameFromDefinition(&definition.Game{
			Players: []*definition.Player{{ID: playerID}},
			Battlefield: &definition.Battlefield{
				Permanents: []*definition.Permanent{
					forestDefinition,
					islandDefinition,
				},
			},
		})
		got := GetAvailableMana(game, playerID, mockApply).ManaString()
		want := "{G}"
		if got != want {
			t.Errorf("GetAvailableMana() = %v; want %v", got, want)
		}
	*/
}
