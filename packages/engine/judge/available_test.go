package judge

import (
	"deckronomicon/packages/game/definition"
	"deckronomicon/packages/game/definition/definitiontest"
	"deckronomicon/packages/state"
	"testing"
)

func TestGetAvailableMana(t *testing.T) {
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
	got := GetAvailableMana(game, playerID).ManaString()
	want := "{G}"
	if got != want {
		t.Errorf("GetAvailableMana() = %v; want %v", got, want)
	}
}
