package judge

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/definition/definitiontest"
	"deckronomicon/packages/state"
	"testing"
)

func TestGetAvailableMana(t *testing.T) {
	// TODO: This tests relies on the game state being updated correctly,
	// by the reducer. Which is not what we should be testing here.

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
	pool, err := GetAvailableMana(game, playerID)
	if err != nil {
		t.Fatalf("GetAvailableMana(); err = %v", err)
	}
	got := pool.ManaString()
	want := "{G}"
	if got != want {
		t.Errorf("GetAvailableMana() = %v; want %v", got, want)
	}
}
